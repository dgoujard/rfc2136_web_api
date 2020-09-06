package rfc2136

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type DnsRecord struct {
	Fqdn string
	Type string
	Values []string
	TTL uint32
	GeneratedId string
}
func (record *DnsRecord ) GenerateId() {
	stringForId := record.Fqdn+record.Type+strings.Join(record.Values,"")+strconv.FormatUint(uint64(record.TTL), 10);
	h := md5.New()
	io.WriteString(h, stringForId)
	record.GeneratedId = fmt.Sprintf("%x", h.Sum(nil))
}
