package context

import (
	"context"
	"net"
	"time"
)

type Value interface {
	Value(key any) any
	Values() map[any]any
	WithValue(key, value any)

	ID() string

	WithUserID(id int)
	UserID() int
	WithUserAgent(userAgent string)
	UserAgent() string
	WithClientIP(IP net.IP)
	ClientIP() net.IP
	WithTraceID(id string)
	TraceID() string
}

const (
	KeyRequestID = "id"
	KeyUserID    = "user_id"
	KeyUserAgent = "user_agent"
	KeyClientIP  = "client_ip"
	KeyTraceID   = "trace_id"
)

var systemKeys = map[string]struct{}{
	KeyRequestID: {},
	KeyUserID:    {},
	KeyUserAgent: {},
	KeyClientIP:  {},
	KeyTraceID:   {},
}

func isSystemKeys(key any) bool {
	strKey, isString := key.(string)
	if !isString {
		return false
	}

	_, exist := systemKeys[strKey]
	return exist
}

func (l *local) Value(key any) any {
	l.mu.RLock()

	value, ok := l.keys[key]
	l.mu.RUnlock()
	if !ok {
		return l.base.Value(key)
	}

	return value
}

func (l *local) Values() map[any]any {
	var values = make(map[any]any, len(l.keys))

	l.mu.RLock()
	for key, value := range l.keys {
		values[key] = value
	}
	l.mu.RUnlock()
	return values
}

func (l *local) WithValue(key, value any) {

	if isSystemKeys(key) {
		// log.Warn("attempting to install a system key. Operation rejected", zap.Any("key", key))
		return // ignore
	}

	l.withValue(key, value)
}

func (l *local) withValue(key, value any) {
	l.base = context.WithValue(l.base, key, value)
	l.mu.Lock()
	l.keys[key] = value
	l.mu.Unlock()
}

func (l *local) ID() string {
	id, _ := l.id()
	return id
}

func (l *local) id() (string, bool) {
	value := l.Value(KeyRequestID)
	id, ok := value.(string)
	return id, ok
}

func (l *local) WithUserID(id int) {
	l.withValue(KeyUserID, id)
}

func (l *local) UserID() int {
	value, ok := l.Value(KeyUserID).(int)
	if !ok {
		return 0
	}
	return value
}

func (l *local) WithUserAgent(userAgent string) {
	l.withValue(KeyUserAgent, userAgent)
}

func (l *local) UserAgent() string {
	value, ok := l.Value(KeyUserAgent).(string)
	if !ok {
		return ""
	}
	return value
}

func (l *local) WithClientIP(userAgent net.IP) {
	l.withValue(KeyClientIP, userAgent)
}

func (l *local) ClientIP() net.IP {
	value, ok := l.Value(KeyClientIP).(net.IP)
	if !ok {
		return net.IP{}
	}
	return value
}

func (l *local) WithTraceID(id string) {
	l.withValue(KeyTraceID, id)
}

func (l *local) TraceID() string {
	value, ok := l.Value(KeyTraceID).(string)
	if !ok {
		return ""
	}
	return value
}

func (l *local) WithTimeout(timeout time.Duration) {
	l.base, l.cancelFunc = context.WithTimeout(l.base, timeout)
}

func (l *local) CopyWithTimeout(timeout time.Duration) Context {
	copyCtx := l.copy()
	copyCtx.base, copyCtx.cancelFunc = context.WithTimeout(copyCtx.base, timeout)
	return copyCtx

}

func (l *local) WithDeadline(d time.Time) {
	l.base, l.cancelFunc = context.WithDeadline(l.base, d)
}

func (l *local) CopyWithDeadline(d time.Time) Context {
	copyCtx := l.copy()
	copyCtx.base, copyCtx.cancelFunc = context.WithDeadline(copyCtx.base, d)
	return copyCtx
}
