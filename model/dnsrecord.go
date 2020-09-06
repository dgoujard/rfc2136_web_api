package model

import (
	"github.com/dgoujard/rfc2136_web_api/util/rfc2136"
	"fmt"
	"github.com/miekg/dns"
	"strings"
)

type DnsRecordForm struct {
	Fqdn         string `form:"required"`
	Type        string `form:"required"`
	Values []string `form:"required"`
	TTL uint32 `form:"required"`
	GeneratedId string
}

func (d *DnsRecordForm) ToModel(zoneName string) (*rfc2136.DnsRecord, error) {
	zoneNameLastChar := zoneName[len(zoneName)-1:]
	var zoneNameWithoutDot string
	if zoneNameLastChar == "."{
		zoneNameWithoutDot = zoneName[:len(zoneName)-1]
	}else {
		zoneNameWithoutDot = zoneName
	}
	if !strings.Contains(d.Fqdn, zoneNameWithoutDot) { //Si pas de zone dans le Fqdn il faut l'ajouter
		FqdnLastChar := d.Fqdn[len(d.Fqdn)-1:]
		if FqdnLastChar != "."{
			d.Fqdn += "."
		}
		d.Fqdn += zoneNameWithoutDot+"."
	}
	d.Fqdn = dns.Fqdn(d.Fqdn)
	fmt.Println(d.Fqdn)
	return &rfc2136.DnsRecord{
		Fqdn:         d.Fqdn,
		Type:        d.Type,
		Values:      d.Values,
		TTL:   d.TTL,
		GeneratedId:   d.GeneratedId,
	}, nil
}