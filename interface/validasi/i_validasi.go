package ivalidasi

import (
	"context"
)

type UseCase interface {
	Barber(ctx context.Context, ID int) error
}
