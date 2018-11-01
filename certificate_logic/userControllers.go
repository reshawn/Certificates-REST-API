package certificates

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getUserCerts(id string) (certCollection, bool) {
	userCerts := certCollection{}
	found := false
	for _, element := range certs {
		if element.OwnerID == id {
			userCerts = append(userCerts, element)
			found = true
		}
	}
	return userCerts, found
}

//return all certificates by the specified user
func userCerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//log.Println(vars)
	id := vars["userID"] // id of user
	log.Println("Get certs of user", id)

	userCerts, found := getUserCerts(id)

	//if user not found
	if !found {
		log.Println("User not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404: User not found"))
		return
	}

	data, _ := json.Marshal(userCerts) //convert data returned to json

	//create and write http response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return

}
