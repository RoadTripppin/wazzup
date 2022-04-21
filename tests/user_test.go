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

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateFailureWithInvalidToken(t *testing.T) {
	var body = []byte(`{"email":"Aditya"}`)

	req, err := http.NewRequest("POST", "/user/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestSearchSuccess(t *testing.T) {
	var body = []byte(`{"querystring": "a"}`)

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.SearchUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSearchFailureWithInvalidToken(t *testing.T) {
	var body = []byte(`{"querystring": "abhiabkbcksb"}`)

	req, err := http.NewRequest("POST", "/user/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.SearchUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestInteractionSuccess(t *testing.T) {
	var body = []byte(``)

	req, err := http.NewRequest("GET", "/user/interacted", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetInteractedUsers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestInteractionFailureWithInvalidToken(t *testing.T) {
	var body = []byte(``)

	req, err := http.NewRequest("GET", "/user/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer yJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA2MjMzMjUsInVzZXJfaWQiOiJlOTkwNTdlZS0yMDU2LTQ2OTMtYjAyZi01YTI2Mzc4OWY5ZDkifQ.taxsx4RXOyEp_x5rJCRmeFdqEbskp-aNaZK6WuVAdNU")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.GetInteractedUsers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
