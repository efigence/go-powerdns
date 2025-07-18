package schema

import "time"

type PDNSDomain struct {
	ID             int      `json:"id"`
	Zone           string   `json:"zone"`
	Masters        []string `json:"masters"`
	NotifiedSerial int      `json:"notified_serial"`
	Serial         int      `json:"serial"`
	LastCheck      int      `json:"last_check"`
	Kind           string   `json:"kind"`
}

func NewPDNSDomainList(domains []DNSDomain) []PDNSDomain {
	d := []PDNSDomain{}
	for i, r := range domains {
		zone := r.Name + "."
		d = append(d, PDNSDomain{
			ID:             i,
			Zone:           zone,
			Masters:        nil,
			NotifiedSerial: int(r.Serial),
			Serial:         int(r.Serial),
			LastCheck:      int(time.Now().Unix()),
			Kind:           "native",
		})
	}
	return d
}
