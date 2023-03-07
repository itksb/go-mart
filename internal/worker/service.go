package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itksb/go-mart/internal/domain"
	"github.com/itksb/go-mart/pkg/logger"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"time"
)

type AccrualWorker struct {
	url    string
	db     *WorkerStorage
	l      logger.Interface
	cancel context.CancelFunc
}

type AccrualResponse struct {
	UserID  string
	Order   string             `json:"order"`
	Status  domain.OrderStatus `json:"status"`
	Accrual float32            `json:"accrual"`
}

func NewAccrualWorker(accrualURL string, db *WorkerStorage, l logger.Interface) (*AccrualWorker, error) {
	return &AccrualWorker{
		url: accrualURL,
		db:  db,
		l:   l,
	}, nil
}

func (a *AccrualWorker) Run(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)
	ctx2, cancel := context.WithCancel(gctx)
	a.cancel = cancel
	group.Go(func() error {
		for loop := true; loop; {
			select {
			case <-ctx2.Done():
				loop = false
			default:
				orders, err := a.db.LoadOrdersToWork(ctx2)
				if err != nil {
					return err
				}
				for _, order := range orders {
					req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/orders/%s", a.url, order.ID), nil)
					if err != nil {
						a.l.Errorf("Worker. http.NewRequest creating error, %s", err)
						continue
					}
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						a.l.Errorf("Worker. http.NewRequest error, %s", err)
						continue
					}
					resBytes, err := io.ReadAll(res.Body)
					if err != nil {
						a.l.Errorf("Worker. %s", err)
						continue
					}
					defer res.Body.Close()

					if res.StatusCode == http.StatusTooManyRequests {
						t, err := time.ParseDuration(res.Header.Get("Retry-After") + "s")
						if err != nil {
							a.l.Errorf("Worker. Parse duration Header error %s", err)
							continue
						}
						time.Sleep(t * time.Second)
						continue
					}
					var accrualResponse AccrualResponse
					if len(resBytes) > 0 {
						if err = json.Unmarshal(resBytes, &accrualResponse); err != nil {
							a.l.Errorf("Worker. Parse json error %s", err)
							continue
						}
					}

					if res.StatusCode == http.StatusOK {
						if accrualResponse.Status == order.Status {
							continue
						} else {

							accrualResponse.UserID = order.UserID
							err := a.db.UpdateOrderStatus(accrualResponse)
							if err != nil {
								a.l.Errorf("Worker. UpdateOrderStatus error %s", err)
								continue
							}

						}
					}
					if res.StatusCode == http.StatusNoContent && accrualResponse.Order != "" {
						err := a.db.UpdateOrderStatus(accrualResponse)
						if err != nil {
							a.l.Errorf("Worker. UpdateOrderStatus error %s", err)
							continue
						}
						continue
					}
				}
				// период опроса БД
				time.Sleep(5 * time.Second)
			}
		}
		return nil
	})

	return group.Wait()
}

func (a *AccrualWorker) Close() error {
	a.cancel()
	return nil
}
