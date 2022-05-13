package repomocks

import "context"

type TxRunner struct {
}

func (m *TxRunner) RunWithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
