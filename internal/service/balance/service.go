package balance

import (
	"context"
	"github.com/itksb/go-mart/internal/domain"
)

type Service struct {
	db domain.BalanceRepositoryInterface
}

type Summary struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func NewBalanceService(db domain.BalanceRepositoryInterface) (*Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) GetBalanceForUser(ctx context.Context, userID string) (*domain.Balance, error) {
	return s.db.FindByUserID(ctx, userID)
}

func (s *Service) GetSummaryForUserID(ctx context.Context, userID string) (*Summary, error) {
	bal, err := s.db.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	withdrawn, err := s.db.SumWithdrawnByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := Summary{
		Current:   bal.Balance,
		Withdrawn: withdrawn,
	}
	return &summary, nil
}
