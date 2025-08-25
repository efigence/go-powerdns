Backend for PowerDNS remote plugin


## PowerDNS config

    launch+=remote
    remote-connection-string=http:url=http://localhost:63636/dns,post,post_json


## DNS record format

```yaml
records: 
    example.com
    ns: [ns1.example.com]
    owner: hostmaster.efigence.com
    expiry: 36000s # default expiry for new records
    autogenerate_ptr: true # use this domain to populate PTR records
    records:
        "*":
            A: ['1.2.3.4']
            MX:
                - value: mx1.example.com
                - value: mx2.example.com
                  prio: 100
                  ttl: 100s
        "www":
            ttl: 60s
            A:
                - '3.4.5.6'
                - '3.4.5.7'           
```
    

