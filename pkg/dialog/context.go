package dialog

import (
	"context"

	"github.com/puppetlabs/relay/pkg/debug"
)

func FromContext(ctx context.Context) Dialog {
	obj := ctx.Value("dialog.dialog")

	if obj == nil {
		debug.Logf("failed to find dialog in context")
		return newNoopDialog()
	}

	if d, ok := obj.(Dialog); ok {
		return d
	} else {
		debug.Logf("object of type %T found in context when looking for Dialog?", obj)
		return newNoopDialog()
	}
}

func NewContext(ctx context.Context, d Dialog) context.Context {
	return context.WithValue(ctx, "dialog.dialog", d)
}
