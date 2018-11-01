package certificates

import (
	"net/http"
	"testing"
)

//router, executeRequest and checkResponseCode are defined in certControllers_test.go
//this file of unit tests can be considered an extension of that and is separated solely
//for the purposes of separating duties and logic

//TestUserCerts test fetching list of user certificates
func TestUserCerts(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/vvg01/certificates", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

}

//TestUserCertsNonExistUser test fetching list from non-existent user
func TestUserCertsNonExistUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/kh01/certificates", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

}
