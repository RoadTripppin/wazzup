package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoadTripppin/wazzup/controllers"
)

func TestUpdateSuccess(t *testing.T) {
	var body = []byte(`{"name":"Aditya"}`)

	req, err := http.NewRequest("POST", "/user/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA3MjAxNDcsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.tNe7bqlpHITw_f_UU0j0SLcDA_e7QzyzUbTVVxauV8o")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
