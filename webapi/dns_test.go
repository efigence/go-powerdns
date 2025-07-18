package webapi

import (
	"bytes"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
func TestDNS(t *testing.T) {
	backend := memdb.New()
	router, err := New(Config{
		Logger:       zaptest.NewLogger(t).Sugar(),
		AccessLogger: zaptest.NewLogger(t).Sugar(),
		ListenAddr:   "0.0.0.0:12345",
		DNSBackend:   backend,
	}, testFS)
	require.NoError(t, err)
	backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	backend.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "example.com",
		Content: "1.2.3.4",
		Ttl:     61,
	})
	//
	w := httptest.NewRecorder()
	b := bytes.Buffer{}
	b.WriteString(`{"method": "getAllDomains", "parameters": {"include_disabled": true}}`)
	req, _ := http.NewRequest(
		"POST",
		"/dns",
		&b,
	)

	router.r.ServeHTTP(w, req)
	//
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"zone":"example.com."`)
	w = httptest.NewRecorder()
	b = bytes.Buffer{}
	b.WriteString(`{"method": "lookup", "parameters": {"local": "0.0.0.0", "qname": "example.com.", "qtype": "ANY", "real-remote": "0.0.0.0/0", "remote": "0.0.0.0", "zone-id": 0}}`)
	req, _ = http.NewRequest(
		"POST",
		"/dns",
		&b,
	)
	router.r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `1.2.3.4`)
}

func BenchmarkWebBackend_Dns(t *testing.B) {
	backend := memdb.New()
	router, err := New(Config{
		Logger:       zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel)).Sugar(),
		AccessLogger: zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel)).Sugar(),
		ListenAddr:   "0.0.0.0:12345",
		DNSBackend:   backend,
	}, testFS)
	require.NoError(t, err)
	backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	backend.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "example.com",
		Content: "1.2.3.4",
		Ttl:     61,
	})
	//
	w := httptest.NewRecorder()

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		b := bytes.Buffer{}
		b.WriteString(`{"method": "getAllDomains", "parameters": {"include_disabled": true}}`)
		req, _ := http.NewRequest(
			"POST",
			"/dns",
			&b,
		)
		router.r.ServeHTTP(w, req)
	}
}
func TestDomainMetadata(t *testing.T) {
	backend := memdb.New()
	router, err := New(Config{
		Logger:       zaptest.NewLogger(t).Sugar(),
		AccessLogger: zaptest.NewLogger(t).Sugar(),
		ListenAddr:   "0.0.0.0:12345",
		DNSBackend:   backend,
	}, testFS)
	require.NoError(t, err)
	backend.AddDomain(schema.DNSDomain{
		Name: "example.com",
		NS:   []string{"ns1.example.com"},
	})
	backend.AddRecord(schema.DNSRecord{
		QType:   "A",
		QName:   "example.com",
		Content: "1.2.3.4",
		Ttl:     61,
	})
	//
	w := httptest.NewRecorder()
	b := bytes.Buffer{}
	b.WriteString(`{"method":"getalldomainmetadata", "parameters":{"name":"example.com"}}`)
	req, _ := http.NewRequest(
		"POST",
		"/dns",
		&b,
	)

	router.r.ServeHTTP(w, req)
	//
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"result":{}`)

}
