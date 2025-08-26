package yamldb

import (
	"fmt"
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/efigence/go-powerdns/backend/yamlloader"
	"github.com/efigence/go-powerdns/schema"
	"github.com/k0kubun/pp"
	"go.uber.org/zap"
	"io/fs"
	"math"
	"math/rand/v2"

	"path/filepath"
	"strings"
	"time"
)

var serialShardsInterval = 864000
var serialShards = math.MaxUint32 / serialShardsInterval

type YAMLDB struct {
	db *memdb.MemDomains
	l  *zap.SugaredLogger
}

func New(log *zap.SugaredLogger) (*YAMLDB, error) {
	backend := YAMLDB{}
	backend.l = log
	backend.db = memdb.New(log.Named("memdb"))
	backend.regenSerial()
	return &backend, nil
}

func (db *YAMLDB) regenSerial() {
	if db.db.SerialBase == 0 {
		db.db.SerialBase = uint32(rand.N(serialShards) * serialShardsInterval)
	} else {
		db.db.SerialBase++
	}

}

// TODO figure out how to handle errornous records without failing whole reload
func (db *YAMLDB) LoadFile(file string) error {
	data, err := yamlloader.Load(file)
	if err != nil {
		return err
	}
	for k1, v1 := range data {
		err := db.db.AddDomain(schema.DNSDomain{
			Name:            k1,
			NS:              v1.NS,
			Owner:           v1.Owner,
			Serial:          uint32(time.Now().Second() / 1000),
			Refresh:         86400,
			Retry:           300,
			Expiry:          864000,
			Nxdomain:        100,
			AutogeneratePTR: v1.AutogeneratePTR,
		})
		if err != nil {
			return err
		}
		for k2, v2 := range v1.Records {
			ttl := v1.DefaultExpiry
			if v2.TTL.Seconds() > 0 {
				ttl = v2.TTL
			}
			for _, z := range v2.A {
				var name string
				if len(k2) == 0 {
					name = k1
				} else {
					name = k2 + "." + k1
				}
				db.db.AddRecord(schema.DNSRecord{
					QType:      "A",
					QName:      name,
					Content:    z.String(),
					Ttl:        int32(ttl.Seconds()),
					DomainId:   0,
					ScopeMask:  "",
					AuthString: "",
				})
				if v1.AutogeneratePTR {
					if _, ok := db.db.Domains[schema.GeneratePTRDomainFromIPv4(z)]; !ok {
						err := db.db.AddDomain(schema.DNSDomain{
							Name:            schema.GeneratePTRDomainFromIPv4(z),
							NS:              v1.NS,
							Owner:           v1.Owner,
							Serial:          uint32(time.Now().Second() / 1000),
							Refresh:         86400,
							Retry:           300,
							Expiry:          864000,
							Nxdomain:        100,
							AutogeneratePTR: false,
						})
						if err != nil {
							pp.Print(err)
						}
					}
					err := db.db.AddRecord(schema.DNSRecord{
						QType:      "PTR",
						QName:      schema.GeneratePTRFromIPv4(z),
						Content:    name,
						Ttl:        int32(ttl.Seconds()),
						DomainId:   0,
						ScopeMask:  "",
						AuthString: "",
					})
					if err != nil {
						pp.Print(err)
					}
				}
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
		er := db.LoadFile(path)
		if er != nil {
			return fmt.Errorf("error while parsing %s: %s", path, er)
		}
		return er
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
	n, _ := New(db.l)
	db.regenSerial()
	n.db.SerialBase = db.db.SerialBase
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

func (db *YAMLDB) DumpAllRecords() map[string]map[string][]schema.DNSRecord {
	return db.db.DomainRecords
}
