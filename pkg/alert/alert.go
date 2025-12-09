package alert

import "context"

type Alert interface {
	Notify(ctx context.Context, errMsg string)
}

var _ Alert = (*AlertNoop)(nil)

type AlertNoop struct{}

func (*AlertNoop) Notify(context.Context, string) {}
