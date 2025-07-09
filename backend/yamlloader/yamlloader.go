package yamlloader

import (
	"github.com/efigence/go-powerdns/schema"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlF struct {
	Zones map[string]schema.Domain `yaml:"zones"`
}

func Load(file string) (d map[string]schema.Domain, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	v := YamlF{}
	err = yaml.NewDecoder(f).Decode(&v)
	return v.Zones, err
}
