package ipredir

func (d *ipredirDomains) AddRedirIp(srcIp string, target string) error {
	var err error
	d.Lock()
	d.redirMap[srcIp] = target
	d.Unlock()
	return err
}

func (d *ipredirDomains) DeleteRedirIp(srcIp string) error {
	var err error
	d.Lock()
	if _, ok := d.redirMap[srcIp]; ok {
		delete(d.redirMap, srcIp)
	}
	return err
}

func (d *ipredirDomains) SetRedirIp(rmap map[string]string) error {
	var err error
	d.Lock()
	d.redirMap = rmap
	d.Unlock()
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
