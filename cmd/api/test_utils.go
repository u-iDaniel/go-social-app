package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/u-iDaniel/go-social-app/internal/auth"
	"github.com/u-iDaniel/go-social-app/internal/ratelimiter"
	"github.com/u-iDaniel/go-social-app/internal/store"
	"github.com/u-iDaniel/go-social-app/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)
	app := &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
	return app
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expectedCode, actualCode int) {
	t.Helper()
	if expectedCode != actualCode {
		t.Errorf("expected status code %d, got %d", expectedCode, actualCode)
	}
}
