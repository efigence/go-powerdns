package webapi

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testFS = os.DirFS("..")

func TestPingRoute(t *testing.T) {

	router, err := New(Config{
		Logger:       zaptest.NewLogger(t).Sugar(),
		AccessLogger: zaptest.NewLogger(t).Sugar(),
		ListenAddr:   "0.0.0.0:12345",
	}, testFS)
	require.NoError(t, err)
	//
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/_status/health", nil)
	router.r.ServeHTTP(w, req)
	//
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "state")
}
