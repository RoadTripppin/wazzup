package tests

import (
	"testing"
	"os"
	"github.com/jinzhu/gorm"
	"github.com/RoadTripppin/wazzup/helpers"
)

/*
	This function is responsible for checking whether the DB connection is successfull or not
*/

func TestDBConn(t *testing.T) {
	db_user := os.Getenv("DB_USER")
	db_name := os.Getenv("DB_NAME")
	db_host := os.Getenv("DB_HOST")
	db_password := os.Getenv("DB_PASSWORD")	
	
	db, err := gorm.Open("postgres", "host=" + db_host + " port=5432 user=" + db_user + " password=" + db_password +" dbname=" + db_name + " sslmode=require")

	if err != nil {
		t.Errorf("Connection to DB Failed")
	}
}

/*
To check if email contains @ symbol or not
*/

func TestRegister(t *testing.T) {

	registerUser := helpers.Register("Max", "max919gmail.com", "12345", "max")

	if registerUser["message"] != "all is fine" {
		t.Errorf("Invalid email")
	}
}

/*
To check if name contains special symbols or not while registering
*/

func TestRegister1(t *testing.T) {

	registerUser := helpers.Register("Max#$", "max919@gmail.com", "12345", "max")

	if registerUser["message"] != "all is fine" {
		t.Errorf("Invalid user name")
	}
}

/*
To check if email contains . symbol or not
*/

func TestRegister2(t *testing.T) {

	loginUser := helpers.Login("max919@gmailcom", "12345max")

	if registerUser["message"] != "all is fine" {
		t.Errorf("Invalid email")
	}
}

func TestGetEntries(t *testing.T) {
	req, err := http.NewRequest("POST", "/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetEntries)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"id":1,"first_name":"Krish","last_name":"Bhanushali","email_address":"krishsb@g.com","phone_number":"0987654321"},{"id":2,"first_name":"xyz","last_name":"pqr","email_address":"xyz@pqr.com","phone_number":"1234567890"},{"id":6,"first_name":"FirstNameSample","last_name":"LastNameSample","email_address":"lr@gmail.com","phone_number":"1111111111"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}