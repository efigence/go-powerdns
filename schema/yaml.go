package schema

import (
	"net"
	"time"
)

type MX struct {
	TTL   time.Duration `yaml:"ttl"`
	Prio  string        `yaml:"prio"`
	Value string        `yaml:"value"`
}

type Record struct {
	TTL   time.Duration     `yaml:"ttl"`
	A     []net.IP          `yaml:"A"`
	AAAA  []net.IP          `yaml:"AAAA"`
	MX    []MX              `yaml:"MX"`
	CNAME []string          `yaml:"CNAME"`
	Extra map[string]string `yaml:"extra"`
}
type Domain struct {
	NS              []string      `yaml:"ns"`
	Owner           string        `yaml:"owner"`
	DefaultExpiry   time.Duration `yaml:"default_expiry"`
	AutogeneratePTR bool          `yaml:"autogenerate_ptr"`
	Records         map[string]Record
}
