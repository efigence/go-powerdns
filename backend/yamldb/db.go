package yamldb

import (
	"fmt"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/backend/yamlloader"
	"github.com/efigence/go-powerdns/schema"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
)

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

func (db *YAMLDB) LoadDir(dir string) error {
	filecount := 0
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		filecount++
		return db.LoadFile(path)
	})
	if err != nil {
		return err
	}
	if filecount > 0 {
		return nil
	} else {
		return fmt.Errorf("zero *.yaml files parsed in dir %s", dir)
	}
}

func (db *YAMLDB) UpdateDir(dir string) error {
	n, _ := New()
	err := n.LoadDir(dir)
	if err != nil {
		return err
	}
	// this is technically wrong, atomic.Pointer should be used but it's such
	// PITA to use that we will just hope golang devs wont fuck up implicit atomic pointer writes
	db.db = n.db
	return err

}

func (db *YAMLDB) Lookup(q schema.QueryLookup) ([]schema.DNSRecord, error) {
	return db.db.Lookup(q)
}
func (db *YAMLDB) List(q schema.QueryList) ([]schema.DNSRecord, error) {
	return db.db.List(q)
}
func (db *YAMLDB) ListDomains(disabled bool) ([]schema.DNSDomain, error) {
	return db.db.ListDomains(disabled)
}

func (db *YAMLDB) GetRootDomainFor(s string) (string, error) {
	return db.db.GetRootDomainFor(s)
}
