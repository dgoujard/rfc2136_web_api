package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgoujard/rfc2136_web_api/model"
	"github.com/dgoujard/rfc2136_web_api/util/rfc2136"
	"github.com/dgoujard/rfc2136_web_api/util/validator"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

// swagger:route GET /api/v1/zones Dns getZones
// Return Zones
// responses:
//  200: SwagZones
func (app *App) HandleZoneList(w http.ResponseWriter, r *http.Request) {
	var b []byte
	b, err := json.Marshal(app.zones)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// swagger:route GET /api/v1/zones/{zoneName} Dns getDnsRecords
// Return DNS records for the zone
// responses:
//  200: body:[]SwagDnsRecord
//  404: DnsRecordRespNotFound
func (app *App) HandleZoneRecordsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zone := ctx.Value("zoneName").(string)

	records , err := app.rfc2136.GetRecordsForZone(zone)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	var b []byte
	b, err = json.Marshal(records)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// swagger:route POST /api/v1/zones/{zoneName} Dns createDnsRecord
// Create DNS record
// responses:
//  200: body:SwagDnsRecord
//  404: DnsRecordRespNotFound
func (app *App) HandleCreateZoneRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zone := ctx.Value("zoneName").(string)

	record, err := app.handleRecordPayload(zone, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	record.GenerateId()

	//Recherche dans les records existant si pas déjà le même enregistrement
	records , err := app.rfc2136.GetRecordsForZone(zone)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	existingRecord := searchAndGetDnsRecordById(record.GeneratedId,&records)
	if existingRecord != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, "Le record existe deja")
		return
	}
	err = app.rfc2136.AddRecordForZone(zone,record)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
		return
	}

	record.GenerateId()
	var b []byte
	b, err = json.Marshal(record)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return
}

// swagger:route GET /api/v1/zones/{zoneName}/{recordId} Dns getDnsRecord
// Return DNS record by Id
// responses:
//  200: body:SwagDnsRecord
//  404: DnsRecordRespNotFound
func (app *App) HandleGetZoneRecordsId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zone := ctx.Value("zoneName").(string)
	recordId := chi.URLParam(r, "recordId")

	records , err := app.rfc2136.GetRecordsForZone(zone)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	record := searchAndGetDnsRecordById(recordId,&records)
	if record == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	var b []byte
	b, err = json.Marshal(record)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
// swagger:route DELETE /api/v1/zones/{zoneName}/{recordId} Dns deleteDnsRecord
// Delete DNS record by Id
// responses:
//  204: description: Deleted successfully
//  404: DnsRecordRespNotFound
func (app *App) HandleDeleteZoneRecordsId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	zone := ctx.Value("zoneName").(string)
	recordId := chi.URLParam(r, "recordId")

	records , err := app.rfc2136.GetRecordsForZone(zone)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	record := searchAndGetDnsRecordById(recordId,&records)
	if record == nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	err = app.rfc2136.DeleteRecordForZone(zone, record)
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	http.Error(w, http.StatusText(204), 204)
}
func (app *App) handleRecordPayload(zone string, body io.Reader) (record *rfc2136.DnsRecord, err error) {
	var recordForm = &model.DnsRecordForm{}

	if err = json.NewDecoder(body).Decode(recordForm); err != nil {
		app.logger.Warn().Err(err).Msg("")
		return nil, errors.New(fmt.Sprintf( `{"error": "%v"}`, appErrFormDecodingFailure))
	}

	if err := app.validator.Struct(recordForm); err != nil {
		app.logger.Warn().Err(err).Msg("")

		resp := validator.ToErrResponse(err)
		if resp == nil {
			return nil, errors.New(fmt.Sprintf( `{"error": "%v"}`, appErrFormErrResponseFailure))
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")

			return nil, errors.New(fmt.Sprintf( `{"error": "%v"}`, appErrJsonCreationFailure))
		}

		return nil, errors.New(string(respBody))
	}
	record, err = recordForm.ToModel(zone)
	if err != nil {
		return
	}
	return
}
func searchAndGetDnsRecordById(searchId string, records *[]rfc2136.DnsRecord) *rfc2136.DnsRecord {
	for _, record := range *records {
		if record.GeneratedId == searchId {
			return &record
		}
	}
	return nil
}