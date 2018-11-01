package certificates

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Data altering functions
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//find certificate to be transferred
//check if authed user is the owner of the certificate
// returns a status code where 1: success, 2: user != owner, 3: cert not found
func addTransferToCert(id string, newTrans transfer, user user) int {
	for i, element := range certs {
		if element.ID == id {
			if element.OwnerID == user.ID { //this user owns the cert
				certs[i].Transfer = newTrans //element != pointer to object in certs
				unacceptedTransfers = append(unacceptedTransfers, &certs[i])
				return 1
			}
			return 2

		}
	}
	return 3
}

//check password of user passed, return user object along with auth status
func authenticate(id string, pass string) (user, bool) {
	for _, u := range users {
		if id == u.ID && pass == u.password {
			return u, true
		}
	}
	return user{}, false

}

//Handler functions
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//create transfer from data sent in request
func createTransfer(w http.ResponseWriter, r *http.Request) {
	var newTrans transfer

	vars := mux.Vars(r)
	//log.Println("Create Transfer for cert ",vars)
	id := vars["id"] // id of certificate to be transferred
	log.Println("Attempt to create Transfer for cert ", id)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Error creating transfer", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating transfer"))
		return
	}

	//unmarshal content of request body as a transfer
	err = json.Unmarshal(body, &newTrans)

	//bad json data
	if err != nil {
		log.Println("Error creating transfer", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating transfer, bad JSON data"))
		return
	}

	log.Println("New transfer:", newTrans)

	//auth user
	ownerID, pass, _ := r.BasicAuth()
	user, valid := authenticate(ownerID, pass)
	if !valid {
		log.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Incorrect user credentials"))
		return
	}

	//execute transfer
	createTransferStatus := addTransferToCert(id, newTrans, user)
	//if user does not own cert
	if createTransferStatus == 2 {
		log.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Unauthorized to create transfer, user does not own certificate"))
		return
	}
	//if cert not found
	if createTransferStatus == 3 {
		log.Println("Certificate not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Certificate not found"))
		return
	}

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	createdData, _ := json.Marshal(newTrans)
	w.Write(createdData)
	return

}

func acceptTransfer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//log.Println(vars)
	id := vars["id"] // id of certificate to be transferred
	log.Println("Attempt to accept Transfer for cert ", id)

	ownerID, pass, _ := r.BasicAuth()
	user, valid := authenticate(ownerID, pass)
	if !valid {
		log.Println("Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401: Incorrect user credentials"))
		return
	}

	//find transfer, ensure receiving user is this user, update status
	for index, element := range unacceptedTransfers {
		if element.ID == id {
			if element.Transfer.To == user.Email {
				unacceptedTransfers[index].Transfer.Status = "Accepted"
				unacceptedTransfers[index].OwnerID = user.ID
				unacceptedTransfers = append(unacceptedTransfers[:index], unacceptedTransfers[index+1:]...)

				//create and write http response
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Transfer Accepted"))
				return
			}
			log.Println("Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401: Transfer not intended for this user"))
			return
		}
	}

	log.Println("Transfer not found")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404: Transfer not found"))
	return
}
