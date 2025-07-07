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

type YAMLDB struct {
	db *memdb.MemDomains
}

func New() (*YAMLDB, error) {
	backend := YAMLDB{}
	backend.db = memdb.New()
	return &backend, nil
}
func (db *YAMLDB) LoadFile(file string) error {
	data, err := yamlloader.Load(file)
	if err != nil {
		return err
	}
	for k1, v1 := range data {
		err := db.db.AddDomain(schema.DNSDomain{
			Name:     k1,
			NS:       v1.NS,
			Owner:    v1.Owner,
			Serial:   uint32(time.Now().Second() / 1000),
			Refresh:  86400,
			Retry:    300,
			Expiry:   864000,
			Nxdomain: 100,
		})
		if err != nil {
			return err
		}
		for k2, v2 := range v1.Records {
			ttl := v1.Expiry
			if v2.TTL.Seconds() > 0 {
				ttl = v2.TTL
			}
			for _, z := range v2.A {
				db.db.AddRecord(schema.DNSRecord{
					QType:      "A",
					QName:      k2 + "." + k1,
					Content:    z.String(),
					Ttl:        int32(ttl.Seconds()),
					DomainId:   0,
					ScopeMask:  "",
					AuthString: "",
				})
			}
		}
	}
	return nil
}

func (db *YAMLDB) Lookup(q api.QueryLookup) (schema.DNSRecordList, error) {
	return db.db.Lookup(q)
}
func (db *YAMLDB) List(q api.QueryList) (schema.DNSRecordList, error) {
	return db.db.List(q)
}

func (db *YAMLDB) GetRootDomainFor(s string) (string, error) {
	return db.db.GetRootDomainFor(s)
}
