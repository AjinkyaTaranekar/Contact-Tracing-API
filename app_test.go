package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
)

func TestGetUser(t *testing.T) {
	var app App
	req, err := http.NewRequest("GET", "/users/5fde5d7e586a36e37a58bb81", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(app.getUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{
        "_id": "5fde5d7e586a36e37a58bb81",
        "name": "User1",
        "dateOfBirth": "31-12-2000",
        "phoneNo": "9995558655",
        "emailAddress": "user1@gmail.com",
        "creationTimeStamp": "2020-12-08T20:07:26.628Z"
    }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	var app App
	var jsonStr = []byte(`{
        "name": "User1",
        "dateOfBirth": "31-12-2000",
        "phoneNo": "9995558655",
        "emailAddress": "user1@gmail.com"
    }`)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(app.createUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{
        "_id": "5fde5d7e586a36e37a58bb81",
        "name": "User1",
        "dateOfBirth": "31-12-2000",
        "phoneNo": "9995558655",
        "emailAddress": "user1@gmail.com",
        "creationTimeStamp": "2020-12-08T20:07:26.628Z"
    }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateContact(t *testing.T) {
	var app App
	var jsonStr = []byte(`{
    	"_idOne":"5f83284980ndfs42",
    	"_idTwo":"5f83284980ndfs46"
    }`)

	req, err := http.NewRequest("POST", "/contact", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(app.contact)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{
		"_id": "5fde5d7e586a36e37a58bb81",
		"_idOne":"5f83284980ndfs42",
    	"_idTwo":"5f83284980ndfs46",	
        "creationTimeStamp": "2020-12-08T20:07:26.628Z"
    }`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestContactTracing(t *testing.T) {

	var app App
	req, err := http.NewRequest("GET", "/contact", nil)
	if err != nil {
		t.Fatal(err)
	}
	q := req.URL.Query()
	q.Add("user", "5fde5d7e586a36e37a58bb81")
	q.Add("infection_timestamp", "2020-12-08T20:07:26.628Z")
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(app.contact)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{
        "_id": "5fde5d7e586a36e37a58bb81",
        "name": "User1",
        "dateOfBirth": "31-12-2000",
        "phoneNo": "9995558655",
        "emailAddress": "user1@gmail.com",
        "creationTimeStamp": "2020-12-08T20:07:26.628Z"
    }]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}