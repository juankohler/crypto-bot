package utils

import (
	"context"
	"time"
)

type detachedContext struct {
	parent context.Context
}

func DetachContext(parent context.Context) context.Context {
	return detachedContext{parent: parent}
}

func (ctx detachedContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx detachedContext) Done() <-chan struct{} {
	return nil
}

func (ctx detachedContext) Err() error {
	return nil
}

func (ctx detachedContext) Value(key interface{}) interface{} {
	return ctx.parent.Value(key)
}
