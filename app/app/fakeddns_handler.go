package app

import (
	"github.com/dgoujard/rfc2136_web_api/util/rfc2136"
	"net/http"
	"strings"
)

// swagger:route GET /nic/update FakeDns FakeDDNSNoIP
// Return None
// responses:
//  200: HealthRespOk
//  404: DnsRecordRespNotFound

func (app *App) HandleFakeDDNSNoIP(w http.ResponseWriter, r *http.Request) {
	myip := r.URL.Query().Get("myip")
	ipRecordType := "A"
	if(strings.Contains(myip, ":")) {
		ipRecordType = "AAAA"
	}
	hostname := r.URL.Query().Get("hostname")

	parts := strings.Split(hostname, ".")
	zone := parts[len(parts)-2] + "." + parts[len(parts)-1]
	recordArr := parts[:len(parts)-2]
	recordFqdn := strings.Join(recordArr,".")+"."+zone+"."
	var found = false
	for _, item := range app.zones {
		if item == zone {
			found = true
		}
	}
	if !found {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	zoneRecords , err := app.rfc2136.GetRecordsForZone(zone)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	var zoneRecordNeedUpdate = true;
	for _, zoneRecord := range zoneRecords {
		if recordFqdn == zoneRecord.Fqdn {
			if zoneRecord.Values[0] == myip && zoneRecord.Type == ipRecordType{
				zoneRecordNeedUpdate = false
			}else{
				app.rfc2136.DeleteRecordForZone(zone,&zoneRecord)
			}
		}
	}
	w.WriteHeader(http.StatusOK)

	if zoneRecordNeedUpdate {
		newDnsRecord := rfc2136.DnsRecord{
			Fqdn:        recordFqdn,
			Type:        ipRecordType,
			Values:      []string{myip},
			TTL:         0,
		}
		app.rfc2136.AddRecordForZone(zone,&newDnsRecord)
		w.Write([]byte("Updated"))
	}else{
		w.Write([]byte("No change"))
	}
}