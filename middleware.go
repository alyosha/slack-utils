package utils

import (
	"context"
	"net/http"
	"time"
)

type customContext struct {
	originalCtx context.Context
}

func (c customContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (c customContext) Done() <-chan struct{}             { return nil }
func (c customContext) Err() error                        { return nil }
func (c customContext) Value(key interface{}) interface{} { return c.originalCtx.Value(key) }

func newContextCopy(ctx context.Context) context.Context {
	return customContext{originalCtx: ctx}
}

// RespondAsync is a middleware that allows the handler to process/respond to
// the initial request without being strictly constrained by Slack's response
// timeout logic. It creates a copy of the context minus the original deadline
// and returns a 200 while letting the next handler finish running in a new goroutine.
func RespondAsync(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxCopy := newContextCopy(r.Context())
		reqCopy := r.Clone(ctxCopy)
		go func() {
			next.ServeHTTP(w, reqCopy)
		}()
	})
}
