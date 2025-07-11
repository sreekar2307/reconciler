package recon

import (
	"context"
	"github.com/sreekar2307/reconciler/cmd"
	"time"
)

func Run(ctx context.Context, deps *cmd.Deps) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := deps.Reconciler.Reconcile(ctx); err != nil {
				return err
			}
		}
	}
}
