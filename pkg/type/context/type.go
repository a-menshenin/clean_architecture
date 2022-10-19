package context

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Context interface {
	context.Context

	WithDeadline(d time.Time)
	CopyWithDeadline(d time.Time) Context

	WithTimeout(timeout time.Duration)
	CopyWithTimeout(timeout time.Duration) Context
	Cancel()

	Copy() Context

	Value
}

type local struct {
	base       context.Context
	cancelFunc context.CancelFunc

	// This mutex protects Keys map.
	mu   sync.RWMutex
	keys map[any]any
}

func (l *local) Deadline() (deadline time.Time, ok bool) {
	return l.base.Deadline()
}

func (l *local) Done() <-chan struct{} {
	return l.base.Done()
}

func (l *local) Err() error {
	return l.base.Err()
}

func (l *local) Cancel() {
	if l.cancelFunc != nil {
		l.cancelFunc()
	}
}

func (l *local) Copy() Context {
	return l.copy()
}

func (l *local) copy() *local {
	l.mu.Lock()
	defer l.mu.Unlock()

	keys := make(map[any]any)
	for k, v := range l.keys {
		keys[k] = v
	}

	return &local{
		base:       l.base,
		cancelFunc: l.cancelFunc,
		mu:         sync.RWMutex{},
		keys:       keys,
	}
}

func (l *local) isEmptyID() bool {
	_, ok := l.id()
	return !ok
}

var cancelFunc = func() {}

func Empty() Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
		keys:       make(map[any]any, 0),
		mu:         sync.RWMutex{},
	}

	ctx.withValue(KeyRequestID, uuid.New().String())

	return ctx
}

func New(option interface{}) Context {
	ctx := &local{
		base:       context.Background(),
		cancelFunc: cancelFunc,
		keys:       make(map[any]any, 0),
		mu:         sync.RWMutex{},
	}

	switch baseCtx := option.(type) {
	case gin.Context:
		ctx.base = baseCtx.Request.Context()

		for key, value := range baseCtx.Keys {
			ctx.withValue(key, value)
		}

		for key := range systemKeys {
			ctx.withValue(key, baseCtx.Value(key))
		}

	case *gin.Context:
		ctx.base = baseCtx.Request.Context()
		for key, value := range baseCtx.Keys {
			ctx.withValue(key, value)
		}

		for key := range systemKeys {
			ctx.withValue(key, baseCtx.Request.Context().Value(key))
		}

	case Context:
		ctx.withValue(KeyRequestID, baseCtx.ID())
	case context.Context:
		ctx.base = baseCtx
		for key := range systemKeys {
			ctx.withValue(key, baseCtx.Value(key))
		}
	}

	if ctx.isEmptyID() {
		ctx.withValue(KeyRequestID, uuid.New().String())
	}
	return ctx
}
