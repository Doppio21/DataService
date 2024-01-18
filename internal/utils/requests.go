package utils

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
)

type RequestFunc func(ctx context.Context) error

func ParallelRequest(ctx context.Context, timeout time.Duration, reqs ...RequestFunc) error {
	res := make([]chan error, len(reqs))

	for i := range reqs {
		res[i] = make(chan error)
		i := i
		go func() {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			res[i] <- reqs[i](ctx)
		}()
	}

	var retErr error
	for _, ch := range res {
		err := <-ch
		if err != nil {
			retErr = multierror.Append(retErr, err)
		}
	}

	return retErr
}
