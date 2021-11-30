package main

import (
	"context"
	_ "embed"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


//go:embed db/ddl.sql
var dbSchemaSetup string

type serverCaller struct {
	url    string
	client *http.Client
}

func (sc *serverCaller) call(req *http.Request) (*http.Response, error) {
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

var sc serverCaller

func TestMain(m *testing.M) {
	db, err := setUpdDB("root", "ss", "simple_auth_test", "multiStatements=true")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("drop database simple_auth_test")
	_, err = db.Exec("create database simple_auth_test")
	_, err = db.Exec("use simple_auth_test")
	_, err = db.Exec(dbSchemaSetup)
	if err != nil {
		panic(err)
	}
	
	mux := setUpMux(db)

	server := httptest.NewServer(mux)
	sc = serverCaller{
		url: server.URL,
		client: server.Client(),
	}
	defer server.Close()
	m.Run()
}

func TestCreateUser(t *testing.T)  {
	reqBody := `{"email": "ss@gmail.com", "password": "111111"}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sc.url + "/user", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	successResp, err := sc.call(req)
	if err != nil {
		t.Fatalf("can't call a server: %v", err)
	}
	defer successResp.Body.Close()
	successBody, err := ioutil.ReadAll(successResp.Body)
	expectedSuccessBody := `{"userId": 1}`
	if successResp.StatusCode != http.StatusCreated {
		t.Errorf("expected status '%d', got '%d'", successResp.StatusCode, http.StatusCreated)
	}
	if string(successBody) != expectedSuccessBody {
		t.Errorf("expected success body '%s', got '%s'", expectedSuccessBody, string(successBody))
	}

	dupEmailResp, err := sc.call(req)
	if err != nil {
		t.Fatalf("can't call a server: %v", err)
	}
	defer dupEmailResp.Body.Close()
	dupEmailBody, err := ioutil.ReadAll(dupEmailResp.Body)
	expectedDupEmailBody := `{"error": 0, "error_desc": "such email is already exists"}`
	if dupEmailResp.StatusCode != http.StatusConflict {
		t.Errorf("expected req status '%d', got '%d'", dupEmailResp.StatusCode, http.StatusConflict)
	}
	if string(dupEmailBody) != expectedDupEmailBody {
		t.Errorf("expected dup email body is '%s', got '%s'", expectedDupEmailBody, string(dupEmailBody))
	}
}

func TestSecurityOnlyAcceptJson(t *testing.T)  {
	reqBody := `{"email": "ss@gmail.com", "password": "111111"}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sc.url + "/user", strings.NewReader(reqBody))
	//req.Header.Set("Content-Type", "application/json") // testng this
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	resp, err := sc.call(req)
	if err != nil {
		t.Fatalf("can't call a server: %v", err)
	}
	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("expected req status '%d', got '%d'", resp.StatusCode, http.StatusUnsupportedMediaType)
	}
}

func TestResponseHeaders(t *testing.T)  {
	reqBody := `{"email": "ss@gmail.com", "password": "111111"}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, sc.url + "/user", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json") // testng this
	if err != nil {
		t.Fatalf("can't create request: %v", err)
	}
	resp, err := sc.call(req)
	if err != nil {
		t.Fatalf("can't call a server: %v", err)
	}

	headers := []struct{
		n  string
		ev string
	}{
		{n: "Content-Type",            ev: "application/json;charset=utf-8"},
		{n: "X-Content-Type-Options",  ev: "nosniff"},
		{n: "X-Frame-Options",         ev: "DENY"},
		{n: "X-XSS-Protection",        ev: "0"},
		{n: "Cache-Control",           ev: "no-store"},
		{n: "Content-Security-Policy", ev: "default-src 'none'; frame-ancestors 'none'; sandbox"},
	}

	for _, header := range headers {
		t.Run(header.n, func(t *testing.T) {
			actualValue := resp.Header.Get(header.n)
			if header.ev != actualValue {
				t.Errorf("expected '%s', got '%s'", header.ev, actualValue)
			}
		})
	}
}
