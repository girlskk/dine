package eventbus

import (
	"context"

	"github.com/gookit/event"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
)

type baseEventContext struct {
	event.BasicEvent
	ds  domain.DataStore
	ctx context.Context
}
