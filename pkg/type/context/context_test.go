package context

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	t.Run("test empty", func(t *testing.T) {
		ctx := Empty()

		assert.NotEmpty(t, ctx)
		assert.NotEqual(t, ctx.ID(), "")
		assert.Equal(t, ctx.UserID(), 0)
		assert.Equal(t, ctx.UserAgent(), "")
		assert.Equal(t, ctx.ClientIP(), net.IP{})
	})

	t.Run("test set system variables with constructor", func(t *testing.T) {
		ctx := Empty()

		ctx.WithUserID(1)
		assert.Equal(t, ctx.UserID(), 1)

		ctx.WithClientIP(net.ParseIP("127.0.0.1"))
		assert.Equal(t, ctx.ClientIP(), net.ParseIP("127.0.0.1"))

		ctx.WithUserAgent("test-user-agent")
		assert.Equal(t, ctx.UserAgent(), "test-user-agent")

		ctx.WithTraceID("test-trace-id")
		assert.Equal(t, ctx.TraceID(), "test-trace-id")
	})

	t.Run("test set system with variables", func(t *testing.T) {
		ctx := Empty()

		for key, value := range systemKeys {
			ctx.WithValue(key, value)
		}
		assert.NotEqual(t, ctx.ID(), "")
		assert.Equal(t, ctx.UserID(), 0)
		assert.Equal(t, ctx.UserAgent(), "")
		assert.Equal(t, ctx.ClientIP(), net.IP{})
	})

	t.Run("test set variables", func(t *testing.T) {
		ctx := Empty()

		ctx.WithValue("test", "test-value")
		assert.Equal(t, ctx.Value("test"), "test-value")
		ctx.WithValue("test-2", "test-2-value")
		assert.Len(t, ctx.Values(), int(3))
	})

	t.Run("test with deadline", func(t *testing.T) {
		ctx := Empty()

		ctx.WithDeadline(time.Now().Add(time.Second))

		for flag := true; flag; {
			select {
			case <-ctx.Done():
				flag = false
			}
		}

		assert.Error(t, ctx.Err(), errors.New("context deadline exceeded"))
	})

	t.Run("test with timeout", func(t *testing.T) {
		ctx := Empty()

		ctx.WithTimeout(time.Second)

		for flag := true; flag; {
			select {
			case <-ctx.Done():
				flag = false
			}
		}

		assert.Error(t, ctx.Err(), errors.New("context deadline exceeded"))
	})

	t.Run("test copy with variables", func(t *testing.T) {
		ctx := Empty()

		ctx.WithValue("test", "test-value")
		ctx.WithValue("test-2", "test-2-value")

		copyCtx := ctx.Copy()

		assert.Equal(t, ctx.Value("test"), copyCtx.Value("test"))
		assert.Equal(t, ctx.Value("test-2"), copyCtx.Value("test-2"))

		copyCtx.WithValue("test", "test-copy-value")
		copyCtx.WithValue("test-3", "test-3-value")

		assert.NotEqual(t, ctx.Value("test-3"), copyCtx.Value("test-3"))
	})

	t.Run("test copy with variables", func(t *testing.T) {
		ctx := Empty()

		ctx.WithValue("test", "test-value")
		ctx.WithValue("test-2", "test-2-value")

		copyCtx := ctx.Copy()

		assert.Equal(t, ctx.Value("test"), copyCtx.Value("test"))
		assert.Equal(t, ctx.Value("test-2"), copyCtx.Value("test-2"))

		copyCtx.WithValue("test", "test-copy-value")
		copyCtx.WithValue("test-3", "test-3-value")

		assert.NotEqual(t, ctx.Value("test-3"), copyCtx.Value("test-3"))
	})

	t.Run("test copy with timeout", func(t *testing.T) {
		ctx := Empty()

		copyCtx := ctx.CopyWithTimeout(time.Second)

		for flag := true; flag; {
			select {
			case <-copyCtx.Done():
				flag = false
			}
		}

		assert.Error(t, copyCtx.Err(), errors.New("context deadline exceeded"))
		assert.NoError(t, ctx.Err())
	})

	t.Run("test copy with deadline", func(t *testing.T) {
		ctx := Empty()

		copyCtx := ctx.CopyWithDeadline(time.Now().Add(time.Second))

		for flag := true; flag; {
			select {
			case <-copyCtx.Done():
				flag = false
			}
		}

		assert.Error(t, copyCtx.Err(), errors.New("context deadline exceeded"))
		assert.NoError(t, ctx.Err())
	})
}
