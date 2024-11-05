package services

import (
	"time"

	"github.com/samber/lo"
)

func ExpiresInToTime(expiresIn *time.Duration) *time.Time {
	if expiresIn == nil {
		return nil
	}

	return lo.ToPtr(time.Now().Add(*expiresIn))
}
