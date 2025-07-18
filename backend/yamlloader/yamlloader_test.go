package yamlloader

import (
	"github.com/efigence/go-powerdns/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	data, err := Load("../../t-data/dns.yaml")
	require.NoError(t, err)
	assert.Equal(t, data, map[string]schema.Domain{
		"example.com": {
			NS:     []string{"ns1.example.com"},
			Expiry: time.Hour * 24,
			Owner:  "hostmaster.example.com",
			Records: map[string]schema.Record{
				"*": {
					A: []net.IP{net.ParseIP("1.2.3.4")},
					MX: []schema.MX{
						{
							Value: "mx1.example.com",
						},
						{
							Value: "mx2.example.com",
							Prio:  "100",
							TTL:   time.Second * 100,
						},
					},
				},
				"www": {
					TTL: time.Second * 3000,
					A: []net.IP{
						net.ParseIP("3.4.5.6"),
						net.ParseIP("3.4.5.7"),
					},
				},
				"": {
					TTL: time.Second * 1234,
					A: []net.IP{
						net.ParseIP("9.9.9.9"),
					},
				},
			},
		},
	})

}
