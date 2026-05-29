package trace

import (
	"context"

	"github.com/dobyte/due/v2/network"
)

const (
	ContextLength    = 25
	connAttrTraceKey = "due.trace.context"
)

type contextKey struct{}

// WithTraceContext stores a neutral trace context carrier in ctx.
func WithTraceContext(ctx context.Context, tc []byte) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(tc) != ContextLength {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, clone(tc))
}

// TraceContextFromContext extracts a neutral trace context carrier from ctx.
func TraceContextFromContext(ctx context.Context) []byte {
	if ctx == nil {
		return nil
	}
	tc, _ := ctx.Value(contextKey{}).([]byte)
	if len(tc) != ContextLength {
		return nil
	}
	return clone(tc)
}

// StoreTraceContext stores a neutral trace context carrier on a network connection.
func StoreTraceContext(conn network.Conn, tc []byte) {
	if conn == nil || len(tc) != ContextLength {
		return
	}
	conn.Attr().Set(connAttrTraceKey, clone(tc))
}

// ClearTraceContext removes a trace context carrier from a network connection.
func ClearTraceContext(conn network.Conn) {
	if conn == nil {
		return
	}
	conn.Attr().Del(connAttrTraceKey)
}

// TraceContextFromConn extracts a neutral trace context carrier from a network connection.
func TraceContextFromConn(conn network.Conn) []byte {
	if conn == nil {
		return nil
	}
	tc, _ := conn.Attr().Get(connAttrTraceKey)
	raw, _ := tc.([]byte)
	if len(raw) != ContextLength {
		return nil
	}
	return clone(raw)
}

func clone(tc []byte) []byte {
	return append([]byte(nil), tc...)
}
