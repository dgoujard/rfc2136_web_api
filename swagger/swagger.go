// API for RFC2136 DNS server
//
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - http_auth: []
//
//     SecurityDefinitions:
//     http_auth:
//          type: basic
//
// swagger:meta
package swagger

import "net/http"

// Success response
// swagger:response HealthRespOk
type swaggHealthResp struct {}

// Not found response
// swagger:response DnsRecordRespNotFound
type swaggDnsRecordNotFound struct {}

// Zone list
// swagger:response SwagZones
type SwagZones struct {
	// in: body
	Body *[]string
}

// DNSRecord response
// swagger:model SwagDnsRecord
type SwagDnsRecord struct {
	Fqdn string
	Type string
	Values []string
	TTL uint32
	GeneratedId string
}

// swagger:parameters createDnsRecord
type AddDnsRecordParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Dnsrecord object that needs to be added to the store
	  In: body
	*/
	Body *SwagDnsRecord
}

// swagger:parameters getDnsRecords getDnsRecord createDnsRecord deleteDnsRecord
type SwagPathZoneName struct {
	// Zone name for the query
	//
	// required: true
	// in: path
	// example: "dgoujard.network"
	ZoneName string `json:"zoneName"`
}
// swagger:parameters getDnsRecord deleteDnsRecord
type SwagPathRecordId struct {
	// Record ID for the query
	//
	// required: true
	// in: path
	RecordId string `json:"recordId"`
}