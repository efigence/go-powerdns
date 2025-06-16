package yamldb

import (
	"github.com/efigence/go-powerdns/api"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/backend/schema"
	"github.com/efigence/go-powerdns/backend/yamlloader"
	"time"
)

func asdf() {
	_, _ = api.New(api.CallbackList{})
}

func New(file string) (api.DomainBackend, error) {
	data, err := yamlloader.Load(file)
	if err != nil {
		return nil, err
	}
	db, _ := memdb.New()
	for k1, v1 := range data {
		db.AddDomain(schema.DNSDomain{
			Name:     k1,
			NS:       v1.NS,
			Owner:    v1.Owner,
			Serial:   uint32(time.Now().Second() / 1000),
			Refresh:  86400,
			Retry:    300,
			Expiry:   864000,
			Nxdomain: 100,
		})
		for k2, v2 := range v1.Records {
			for _, z := range v2.A {
				db.AddRecord(schema.DNSRecord{
					QType:      "A",
					QName:      k2 + "." + "k1",
					Content:    z.String(),
					Ttl:        int32(v1.Expiry.Seconds()),
					DomainId:   0,
					ScopeMask:  "",
					AuthString: "",
				})
			}

		}

	}
	return db, err
}
