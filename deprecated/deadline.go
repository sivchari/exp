package deprecated

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var (
	once sync.Once
	t    time.Time
)

func Setoff(tt time.Time) {
	if tt.Before(time.Now()) {
		panic("time is in the past")
	}
	once.Do(func() {
		t = tt
	})
}

func SetoffDuration(d time.Duration) {
	Setoff(time.Now().Add(d))
}

type DeprecatedDeadlineRoundTrip struct{}

var _ http.RoundTripper = (*DeprecatedDeadlineRoundTrip)(nil)

func (d *DeprecatedDeadlineRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.IsZero() {
		return nil, fmt.Errorf("deadline not set")
	}
	if time.Now().After(t) {
		slog.Warn("deadline exceeded", slog.String("reason", fmt.Sprintf("the time setoff is %s", t.String())))
	}
	return http.DefaultTransport.RoundTrip(req)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if t.IsZero() {
			http.Error(w, "deadline not set", http.StatusInternalServerError)
		}
		if time.Now().After(t) {
			slog.Warn("deadline exceeded", slog.String("reason", fmt.Sprintf("the time setoff is %s", t.String())))
		}
		h.ServeHTTP(w, r)
	})
}
