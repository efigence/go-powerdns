package schema

import (
	"net"
	"strings"
)

// generate array of domains from subdomain, specific -> generic
func ExpandDNSName(name string) ([]string, error) {
	var s []string
	var err error

	parts := strings.Split(name, `.`)
	for i := 0; i < len(parts); i++ {
		s = append(s, strings.Join(parts[i:], `.`))
	}
	return s, err
}

func GeneratePTRFromIPv4(ip net.IP) string {
	octets := strings.Split(ip.String(), ".")
	return octets[3] + "." + octets[2] + "." + octets[1] + "." + octets[0] + "." + "in-addr.arpa"
}

func GeneratePTRDomainFromIPv4(ip net.IP) string {
	octets := strings.Split(ip.String(), ".")
	return octets[2] + "." + octets[1] + "." + octets[0] + "." + "in-addr.arpa"
}
