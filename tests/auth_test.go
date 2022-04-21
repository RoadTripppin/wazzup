package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoadTripppin/wazzup/controllers"
)

func TestRegisterSuccess(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"adi@hotmail.com", "password": "98765", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header("Authorization" : "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE2NTA0NTI3MDgsInVzZXJfaWQiOiI2NTY2YTlhOS01YTA2LTRhMzQtYTY1Yi02MWIxOWRjZGFiZTgifQ.zg3il-ep3wbLwIL3xudzsyayr-J9VQtmklGpYnJHuYw")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestRegisterFailureWithInvalidEmail(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"adihotmail.com", "password": "98765", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithInvalidPassword(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"adi@hotmail.com", "password": "965", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithEmptyEmail(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"", "password": "98765", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithEmptyPassword(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"adi@hotmail.com", "password": "", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithEmptyEmailAndPassword(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"", "password": "", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithNoEmail(t *testing.T) {
	var body = []byte(`{"name":"Adi", "password": "98765", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithNoPassword(t *testing.T) {
	var body = []byte(`{"name":"Adi", "email":"adi@hotmail.com", "profilepic": ""}`)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestRegisterFailureWithNoBody(t *testing.T) {
	// var body = []byte(`{}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Register)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
}

func TestLoginSuccess(t *testing.T) {
	var body = []byte(`{"email":"adi@hotmail.com","password": "98765"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLoginFailureWithInvalidEmail(t *testing.T) {
	var body = []byte(`{"email":"adihotmail.com","password": "98765"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithInvalidPassword(t *testing.T) {
	var body = []byte(`{"email":"adi@hotmail.com","password": "9865"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithEmptyEmail(t *testing.T) {
	var body = []byte(`{"email":"","password": "98765"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithEmptyPassword(t *testing.T) {
	var body = []byte(`{"email":"adi@hotmail.com","password": ""}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithEmptyEmailAndPassword(t *testing.T) {
	var body = []byte(`{"email":"","password": ""}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithNoEmail(t *testing.T) {
	var body = []byte(`{"password": ""}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithNoPassword(t *testing.T) {
	var body = []byte(`{"email": "adi@hotmail.com"}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func TestLoginFailureWithNoBody(t *testing.T) {
	// var body = []byte(`{}`)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.Login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
