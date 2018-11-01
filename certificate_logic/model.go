package certificates

import "time"

type transfer struct {
	To     string // user email
	Status string
}

type user struct {
	ID       string
	Email    string
	Name     string
	password string
}

type certificate struct {
	ID        string
	Title     string
	CreatedAt time.Time //date of creation
	OwnerID   string
	Year      int
	Note      string
	Transfer  transfer //representing current state of transfer
}
type exception struct {
	Message string
}

type certCollection []certificate
type userCollection []user

/* static data for use in this project as opposed to db.
Package scoped variables are bad practice but for the purposes of
this task they allow for more convenient separation of duties */
var certs = certCollection{
	{
		ID:        "c001",
		Title:     "The Starry Night",
		CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		OwnerID:   "rr01",
		Year:      1889,
		Note:      "",
		Transfer:  transfer{To: "", Status: ""},
	},
	{
		ID:        "c002",
		Title:     "Café Terrace at Night",
		CreatedAt: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
		OwnerID:   "vvg01",
		Year:      1888,
		Note:      "",
		Transfer:  transfer{To: "", Status: ""},
	},
}

var users = userCollection{
	{
		ID:       "rr01",
		Email:    "reshawnramjattan@gmail.com",
		Name:     "Reshawn",
		password: "rrejh3294",
	},
	{
		ID:       "vvg01",
		Email:    "vvg@gmail.com",
		Name:     "Vincent Van Golang",
		password: "vwh39043f",
	},
}

//An array of pointers to certificate objects stored in certs that acts as
//a subset of certificates containing transfers that have not yet been handled
//This makes it easier to search for certificates when processing transfers
var unacceptedTransfers = []*certificate{}
