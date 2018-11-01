package certificates

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Route is meant to be a cleaner way to define and store routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	// list all certificates
	Route{
		"all_certificates",
		"GET",
		"/certificates",
		getAllCerts,
	},
	// create new certificate (with ownerd in header)
	Route{
		"create_certificate",
		"POST",
		"/certificates/create",
		createCert,
	},
	//Get certificate by id
	Route{
		"get_certificate",
		"GET",
		"/certificates/{id}",
		getCert,
	},
	//update existing certificate
	Route{
		"update_certificate",
		"PUT",
		"/certificates/update",
		updateCert,
	},
	//Delete product by id
	Route{
		"delete_certificate",
		"DELETE",
		"/certificates/{id}/delete",
		deleteCert,
	},
	//View certificates belonging to a user
	Route{
		"user_certificates",
		"GET",
		"/users/{userID}/certificates",
		userCerts,
	},
	//Create certificate transfer
	Route{
		"create_transfer",
		"POST",
		"/certificates/{id}/transfers/create",
		createTransfer,
	},
	//Accept certificate transfer
	Route{
		"accept_transfer",
		"PUT",
		"/certificates/{id}/transfers/accept",
		acceptTransfer,
	},
}

//NewRouter Configures a new router to the API based on all above routes
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println("Route: ", route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
