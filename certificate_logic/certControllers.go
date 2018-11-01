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

//search data for certificate by id; this is called in getCert
func lookupCert(id string) (certificate, bool) {
	for _, element := range certs {
		if element.ID == id {
			return element, true
		}
	}
	return certificate{}, false
}

//update cert collection given updated cert, returns a status code
//code 1: update successful
//code 2: update failed due to owner change attempt
//code 3: update failed cert not found
func updateCertCollection(uc certificate) int {
	for index, element := range certs {
		if element.ID == uc.ID {
			if element.OwnerID != uc.OwnerID {
				return 2
			}
			certs[index] = uc
			return 1
		}
	}
	return 3
}

//delete cert from collection given cert id, returns success bool
func deleteCertFromCollection(id string) bool {
	for index, element := range certs {
		if element.ID == id {
			//remove element at index; linear time, can be faster if maintaining order doesn't matter
			certs = append(certs[:index], certs[index+1:]...)
			return true
		}
	}
	return false
}

//Handler functions
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

//get all certificates
func getAllCerts(w http.ResponseWriter, r *http.Request) {
	log.Println("Printing All Certificates")
	data, _ := json.Marshal(certs)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return
}

//get certificate by id
func getCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//log.Println(vars)
	id := vars["id"] // id of certificate to be retrieved
	log.Println("Get cert", id)

	cert, found := lookupCert(id)

	//if cert not found
	if !found {
		log.Println("Certificate not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Certificate not found"))
		return
	}

	data, _ := json.Marshal(cert) //convert data returned to json

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return

}

//create certificate from data sent in request
func createCert(w http.ResponseWriter, r *http.Request) {
	var newCert certificate
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Error creating certificate", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating certificate"))
		return
	}

	//unmarshal content of request body as a certificate
	err = json.Unmarshal(body, &newCert)

	//bad json data
	if err != nil {
		log.Println("Error creating certificate", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error creating certificate"))
		return
	}

	log.Println("New certificate create")

	//owner of cert passed as a custom header
	owner := r.Header.Get("OwnerID")

	// ownerid header missing
	if owner == "" {
		log.Println("Header ownerID required to create certificate")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden: Header ownerID required to create certificate"))
		return
	}
	newCert.OwnerID = owner
	certs = append(certs, newCert)

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	createdData, _ := json.Marshal(newCert)
	w.Write(createdData)
	return

}

//update certificate
func updateCert(w http.ResponseWriter, r *http.Request) {
	var updatedCert certificate
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Error updating certificate", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error updating certificate"))
		return
	}

	//unmarshal content of request body as a certificate
	err = json.Unmarshal(body, &updatedCert)

	//bad json data
	if err != nil {
		log.Println("Error updating certificate", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error updating certificate"))
		return
	}

	log.Println("Updated certificate")

	//changing owner is done by a transfer
	//ownerid header not required for update UNLESS the record must be created

	updateStatus := updateCertCollection(updatedCert)

	//error owner change is attempted
	if updateStatus == 2 {
		log.Println("Owner change attempted; Must be done by transfer")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Owner change attempted; Must be done by transfer"))
		return
	}

	//not found, cert must be created
	if updateStatus == 3 {
		owner := r.Header.Get("OwnerID")

		// ownerid header missing
		if owner == "" {
			log.Println("Header ownerID required to create certificate")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden: Record not found, thus header ownerID required to create certificate"))
			return
		}
		updatedCert.OwnerID = owner
		certs = append(certs, updatedCert)
	}

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	updatedData, _ := json.Marshal(updatedCert)
	w.Write(updatedData)
	return

}

//delete certificate
func deleteCert(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//log.Println(vars)
	id := vars["id"] // id of certificate to be deleted
	log.Println("Attempt to delete cert", id)

	found := deleteCertFromCollection(id)

	//if cert not found
	if !found {
		log.Println("Certificate not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: Certificate not found"))
		return
	}

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 Delete Successful"))
	return

}
