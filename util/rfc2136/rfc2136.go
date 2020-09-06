package rfc2136

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	// Map of supported TSIG algorithms
	tsigAlgs = map[string]string{
		"hmac-md5":    dns.HmacMD5,
		"hmac-sha1":   dns.HmacSHA1,
		"hmac-sha256": dns.HmacSHA256,
		"hmac-sha512": dns.HmacSHA512,
	}
)

type Rfc2136 struct {
	nameserver string
	insecure bool
	tsigSecret string
	tsigKeyname string
	tsigSecretAlg string
}

func New(Host string, Port int, insecure bool, TsigSecret string, TsigSecretAlg string, TsigKeyName string) (*Rfc2136, error) {

	secretAlgChecked, ok := tsigAlgs[TsigSecretAlg]
	if !ok && !insecure {
		return nil, errors.New("Algorithme didn't exist")
	}

	return &Rfc2136{
		nameserver: net.JoinHostPort(Host, strconv.Itoa(Port)),
		insecure:insecure,
		tsigSecret: TsigSecret,
		tsigKeyname: TsigKeyName,
		tsigSecretAlg: secretAlgChecked,
	}, nil
}

func (rfc *Rfc2136 ) GetRecordsForZone(zoneName string) (records []DnsRecord, err error)  {
	zoneName = dns.Fqdn(zoneName)

	t := new(dns.Transfer)
	msg := dns.Msg{}
	msg.SetAxfr(zoneName)
	if !rfc.insecure {
		msg.SetTsig(dns.Fqdn(rfc.tsigKeyname), rfc.tsigSecretAlg, 300, time.Now().Unix())
		t.TsigSecret = map[string]string{dns.Fqdn(rfc.tsigKeyname): rfc.tsigSecret}
	}
	channel, err := t.In(&msg, rfc.nameserver)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}
	remoterecords := make([]dns.RR, 0)

	for e := range channel {
		if e.Error != nil {
			if e.Error == dns.ErrSoa {
				fmt.Println("AXFR error: unexpected response received from the server")
			} else {
				fmt.Println("AXFR error: %v", e.Error)
			}
			continue
		}
		remoterecords = append(remoterecords, e.RR...)
	}
	
	for _, rr := range remoterecords {
		if rr.Header().Class != dns.ClassINET {
			continue
		}
		rrFqdn := rr.Header().Name
		rrTTL := rr.Header().Ttl

		var rrType string
		var rrValues []string
		switch rr.Header().Rrtype {
		case dns.TypeCNAME:
			rrValues = []string{rr.(*dns.CNAME).Target}
			rrType = "CNAME"
		case dns.TypeA:
			rrValues = []string{rr.(*dns.A).A.String()}
			rrType = "A"
		case dns.TypeAAAA:
			rrValues = []string{rr.(*dns.AAAA).AAAA.String()}
			rrType = "AAAA"
		case dns.TypeTXT:
			rrValues = (rr.(*dns.TXT).Txt)
			rrType = "TXT"
		case dns.TypeNS:
			rrValues = []string{(rr.(*dns.NS).Ns)}
			rrType = "NS"
		case dns.TypeMX:
			rrValues = []string{(strconv.Itoa(int(rr.(*dns.MX).Preference))) , (rr.(*dns.MX).Mx)}
			rrType = "MX"
		default:
			fmt.Println("Ignored %s",rr)
			continue // Unhandled record type
		}
		dnsRecord := DnsRecord{
			Fqdn:   rrFqdn,
			Type:   rrType,
			Values: rrValues,
			TTL:    rrTTL,
		}
		dnsRecord.GenerateId()
		records = append(records, dnsRecord)
	}
	return
}
func (rfc *Rfc2136 ) AddRecordForZone(zoneName string, record *DnsRecord) error {
	zoneName = dns.Fqdn(zoneName)

	msg := dns.Msg{}
	msg.SetUpdate(zoneName)

	c := new(dns.Client)
	c.SingleInflight = true

	minTTL := time.Duration(record.TTL) * time.Second
	var ttl = int64(minTTL.Seconds())
	newRR := fmt.Sprintf("%s %d %s %s", dns.Fqdn(record.Fqdn), ttl, record.Type, strings.Join(record.Values, " "))
	rr, err := dns.NewRR(newRR)
	if err != nil {
		return err
	}
	msg.Insert([]dns.RR{rr})
	if !rfc.insecure {
		msg.SetTsig(dns.Fqdn(rfc.tsigKeyname), rfc.tsigSecretAlg, 300, time.Now().Unix())
		c.TsigSecret = map[string]string{dns.Fqdn(rfc.tsigKeyname): rfc.tsigSecret}
	}
	resp, _, err := c.Exchange(&msg, rfc.nameserver)
	if err != nil {
		return err
	}
	if resp != nil && resp.Rcode != dns.RcodeSuccess {
		/*fmt.Printf("Bad dns.Client.Exchange response: %s", resp)
		fmt.Printf("bad return code: %s", dns.RcodeToString[resp.Rcode])*/
		return errors.New(fmt.Sprintf("bad return code: %s", dns.RcodeToString[resp.Rcode]))
	}
	fmt.Println("ok")
	return nil
}
func (rfc *Rfc2136 ) DeleteRecordForZone(zoneName string, record *DnsRecord) error {
	zoneName = dns.Fqdn(zoneName)

	msg := dns.Msg{}
	msg.SetUpdate(zoneName)

	c := new(dns.Client)
	c.SingleInflight = true

	minTTL := time.Duration(record.TTL) * time.Second
	var ttl = int64(minTTL.Seconds())
	newRR := fmt.Sprintf("%s %d %s %s", dns.Fqdn(record.Fqdn), ttl, record.Type, strings.Join(record.Values, " "))
	rr, err := dns.NewRR(newRR)
	if err != nil {
		return err
	}
	msg.Remove([]dns.RR{rr})
	if !rfc.insecure {
		msg.SetTsig(dns.Fqdn(rfc.tsigKeyname), rfc.tsigSecretAlg, 300, time.Now().Unix())
		c.TsigSecret = map[string]string{dns.Fqdn(rfc.tsigKeyname): rfc.tsigSecret}
	}
	resp, _, err := c.Exchange(&msg, rfc.nameserver)
	if err != nil {
		return err
	}
	if resp != nil && resp.Rcode != dns.RcodeSuccess {
		/*fmt.Printf("Bad dns.Client.Exchange response: %s", resp)
		fmt.Printf("bad return code: %s", dns.RcodeToString[resp.Rcode])*/
		return errors.New(fmt.Sprintf("bad return code: %s", dns.RcodeToString[resp.Rcode]))
	}
	return nil
}