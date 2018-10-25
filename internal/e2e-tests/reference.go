package e2etests

import "context"

type Retrier interface {
	Retry(ctx context.Context, command Command) error
}

type Command func() error
