package ipredir

func (d *ipredirDomains) RedirIp(srcIp string, target string) error {
	var err error
	d.redirMap[srcIp] = target
	return err
}
