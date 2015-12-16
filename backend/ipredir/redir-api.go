package ipredir

import (
	"net"
	"errors"
	"fmt"
)


func (d *ipredirDomains) AddRedirIp(srcIp string, target string) error {
	var err error
	src := net.ParseIP(srcIp)
	dst := net.ParseIP(target)
	if src == nil {
		return errors.New(fmt.Sprintf("Not a valid IP: %+v", srcIp))
	}
	if dst == nil {
		return errors.New(fmt.Sprintf("Not a valid IP: %+v", srcIp))
	}
	d.Lock()
	defer 	d.Unlock()
	d.redirMap[src.String()] = dst.String()
	return err
}

func (d *ipredirDomains) DeleteRedirIp(srcIp string) error {
	var err error
	d.Lock()
	defer d.Unlock()
	delete(d.redirMap, srcIp)
	return err
}

func (d *ipredirDomains) SetRedirIp(rmap map[string]string) error {
	var err error
	d.Lock()
	defer d.Unlock()
	d.redirMap = rmap

	return err
}

func (d *ipredirDomains) ListRedirIp() (map[string]string, error) {
	var err error
	res := make(map[string]string)
	d.RLock()
	for k, v := range d.redirMap {
		res[k] = v
	}
	d.RUnlock()
	return res, err
}
