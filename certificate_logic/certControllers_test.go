package certificates

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router = NewRouter()

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func checkResponseCode(t *testing.T, expected, got int) {
	if expected != got {
		t.Errorf("Expected response code %d. Got %d\n", expected, got)
	}
}

// TestGetAllCerts simply tests the response of get all certificates
func TestGetAllCerts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/certificates", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

// TestGetCert simply tests the response of fetching an existing certificate
func TestGetCert(t *testing.T) {
	req, _ := http.NewRequest("GET", "/certificates/c001", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["ID"] != "c001" {
		t.Errorf("Expected certificate to be 'c001'. Got '%v'", m["ID"])
	}
}

// TestGetCert tests fetching a nonexistent certificate
func TestGet404Cert(t *testing.T) {
	req, _ := http.NewRequest("GET", "/certificates/c-1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

}

//TestCreateCert tests the creation of a certificate
func TestCreatCert(t *testing.T) {
	testcert := []byte(`{"ID": "c003","Title": "The Yellow House","CreatedAt": "2009-11-17T20:34:58.651387237Z","OwnerID": "rr01","Year": 1888,"Note": "","Transfer": {"To": "","Status": ""}}`)

	req, _ := http.NewRequest("POST", "/certificates/create", bytes.NewBuffer(testcert))
	req.Header.Set("OwnerID", "rr01")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["ID"] != "c003" {
		t.Errorf("Expected certificate ID to be 'c003'. Got '%v'", m["ID"])
	}

	if m["Title"] != "The Yellow House" {
		t.Errorf("Expected certificate title to be 'The Yellow House'. Got '%v'", m["Title"])
	}

	if m["CreatedAt"] != "2009-11-17T20:34:58.651387237Z" {
		t.Errorf("Expected created at to be '2009-11-17T20:34:58.651387237Z'. Got '%v'", m["CreatedAt"])
	}

}

//TestCreateCertNoOwner tests sending a create request without the required owner header
func TestCreateCertNoOwner(t *testing.T) {
	testcert := []byte(`{"ID": "c003","Title": "The Yellow House","CreatedAt": "2009-11-17T20:34:58.651387237Z","OwnerID": "rr01","Year": 1888,"Note": "","Transfer": {"To": "","Status": ""}}`)

	req, _ := http.NewRequest("POST", "/certificates/create", bytes.NewBuffer(testcert))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)
}

//TestUpdateCert tests a correct update of a certificate
func TestUpdateCert(t *testing.T) {
	testcert := []byte(`{"ID": "c001","Title": "THE YELLOW HOUSE","CreatedAt": "2009-11-17T20:34:58.651387237Z","OwnerID": "rr01","Year": 1888,"Note": "","Transfer": {"To": "","Status": ""}}`)

	req, _ := http.NewRequest("PUT", "/certificates/update", bytes.NewBuffer(testcert))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Title"] != certs[0].Title {
		t.Errorf("Expected certificate title to be 'THE YELLOW HOUSE'. Got '%v'", m["Title"])
	}
}

//TestUpdateCertNew update a cert that does not exist, thereby creating a new one
func TestUpdateCertNew(t *testing.T) {
	testcert := []byte(`{"ID": "c003","Title": "THE YELLOW HOUSE","CreatedAt": "2009-11-17T20:34:58.651387237Z","OwnerID": "rr01","Year": 1888,"Note": "","Transfer": {"To": "","Status": ""}}`)

	req, _ := http.NewRequest("PUT", "/certificates/update", bytes.NewBuffer(testcert))
	req.Header.Set("OwnerID", "rr01")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Title"] != certs[2].Title {
		t.Errorf("Expected certificate title to be 'THE YELLOW HOUSE'. Got '%v'", m["Title"])
	}
}

//TestUpdateCertNewNoOwner update a cert that does not exist, thereby attempting to create new one but omit owner ID header
func TestUpdateCertNewNoOwner(t *testing.T) {
	testcert := []byte(`{"ID": "c004","Title": "Wheatfield with Crows","CreatedAt": "2009-11-17T20:34:58.651387237Z","OwnerID": "rr01","Year": 1890,"Note": "","Transfer": {"To": "","Status": ""}}`)

	req, _ := http.NewRequest("PUT", "/certificates/update", bytes.NewBuffer(testcert))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusForbidden, response.Code)

}

//TestDeleteCert test delete
func TestDeleteCert(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/certificates/c003/delete", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

//TestDeleteNonExistingCert test deleting a cert not in storage
func TestDeleteNonExistingCert(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/certificates/c-/delete", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
