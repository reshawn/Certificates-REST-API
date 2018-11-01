package certificates

import (
	"bytes"
	"net/http"
	"testing"
)

//router, executeRequest and checkResponseCode are defined in certControllers_test.go
//this file of unit tests can be considered an extension of that and is separated solely
//for the purposes of separating duties and logic

//TestCreateTransfer test correctly creating a transfer of cert
func TestCreateTransfer(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c001/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("rr01", "rrejh3294")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	if certs[0].Transfer.Status != "pending" {
		t.Errorf("Expected certificate transfer to be pending. Got '%v'", certs[0].Transfer.Status)
	}

}

//TestCreateTransferBadCredentials test create transfer with incorrect basic auth credentials
func TestCreateTransferBadCredentials(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c001/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("rr01", "rr")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

//TestCreateTransferNonExist test create transfer of non-existent certificate
func TestCreateTransferNonExist(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c-/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("rr01", "rrejh3294")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

}

//TestCreateTransferNotMine test create transfer of certificate not owned by auth user
func TestCreateTransferNotMine(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c001/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("vvg01", "vwh39043f")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

//TestAcceptTransfer test correctly accepting a transfer of cert
func TestAcceptTransfer(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c001/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("rr01", "rrejh3294")
	response := executeRequest(req)

	req, _ = http.NewRequest("PUT", "/certificates/c001/transfers/accept", nil)
	req.SetBasicAuth("vvg01", "vwh39043f")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if certs[0].OwnerID != "vvg01" {
		t.Errorf("Expected certificate owner to be changed. Got '%v'", certs[0].OwnerID)
	}

}

//TestAcceptTransferBadCredentials test attempting to accept transfer with incorrect credentials
func TestAcceptTransferBadCredentials(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c001/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("rr01", "rrejh3294")
	response := executeRequest(req)

	req, _ = http.NewRequest("PUT", "/certificates/c001/transfers/accept", nil)
	req.SetBasicAuth("vvg01", "vvg")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

}

//TestAcceptTransferNonExist test attempting to accept transfer that has not been created
func TestAcceptTransferNonExist(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/certificates/c002/transfers/accept", nil)
	req.SetBasicAuth("vvg01", "vwh39043f")
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

}

//TestAcceptTransferNotYours test attempting to accept transfer not intended for auth user
func TestAcceptTransferNotYours(t *testing.T) {
	testtrans := []byte(`{"To": "vvg@gmail.com","Status": "pending"}`)

	req, _ := http.NewRequest("POST", "/certificates/c002/transfers/create", bytes.NewBuffer(testtrans))
	req.SetBasicAuth("vvg01", "vwh39043f")
	response := executeRequest(req)

	req, _ = http.NewRequest("PUT", "/certificates/c002/transfers/accept", nil)
	req.SetBasicAuth("rr01", "rrejh3294")
	response = executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	if body := response.Body.String(); body != "401: Transfer not intended for this user" {
		t.Errorf("Expected transfer not for this user error. Got %s", body)
	}

}
