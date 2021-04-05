package usecase

import (
	"context"
)

type TpbankUsecase interface {
	sayHello(ctx context.Context) error
}
