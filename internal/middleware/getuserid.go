package middleware

import "context"

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxUser).(string)
	if !ok {
		return "", false
	}
	return userID, true
}
