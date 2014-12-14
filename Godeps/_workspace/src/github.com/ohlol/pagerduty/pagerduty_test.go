package pagerduty

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type pagerDutyIncidentResponseMock struct {
	request *http.Request
}

func (m *pagerDutyIncidentResponseMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.request = r
	w.Write([]byte(`{"incidents": [{"id":"1234","incident_number":1,"created_on":"date_here","status":"resolved","html_url":"url_here","incident_key":"incident_key_here","service":{},"assigned_to_user":null,"trigger_summary_data":{},"trigger_details_html_url":"url_here","last_status_change_on":"date_here","last_status_change_by":null,"number_of_escalations":0}],"limit":100,"offset":0,"total":1}`))
}

func TestIncidentsResponse(t *testing.T) {
	handler := &pagerDutyIncidentResponseMock{}
	server := httptest.NewServer(handler)
	account := SetupAccount("testing", "1234")
	account.url = server.URL

	defer server.Close()

	incidents, err := account.Incidents(nil)
	if err != nil {
		t.Error(err)
	}

	if len(incidents) != 1 {
		t.Error("Did not get back 1 incident")
	}
}
