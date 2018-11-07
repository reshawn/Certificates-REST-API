

# Certificates REST API
A few simple RESTful endpoints designed as if they were to be consumed by the Verisart front-end app.

- Perform basic CRUD (create, read, update, delete) operations on JSON certificate objects stored in the back-end. (static data for the purposes of this task)
- View certificates owned by a user
- Transfer certificate ownership from one user to another, supported by basic authorization.

The documentation below walks through setting up and running the code as well as detailed examples of how to consume the endpoints.

## Setup
### Golang
follow the official [docs](https://golang.org/doc/install) for installation and setup of the development environment
### Project
1. After installing Golang, a GOPATH environment variable should be set.
In the directory "[GOPATH]/src/" unzip these files.
2. Install the project dependencies  
`// These libraries are for routing and handling requests` <br>
`$ go get "github.com/gorilla/mux"`<br>
`$ go get "github.com/gorilla/context"`<br>
`$ go get "github.com/gorilla/handlers"`<br>
3. To run the project locally, run `go run main.go`
4. The API is ready to be opened now! Go to http://localhost:8080 in your browser
(the port can be edited in the main.go file)

### Docker
To run the [published](https://hub.docker.com/r/reshawn/certificates/) docker container, run
` docker run --rm -it -p 8080:8080 reshawn/certificates`

### Running Unit Tests
Change directory to the certificate_logic folder and run `go test`


**Note**: A JSON formatter extension or API development environment like [Postman](https://www.getpostman.com/apps) will be useful for trying out the endpoint consumption



## API Documentation & Usage
### 1. View All Certificates
- **Endpoint Name** - `all_certificates`    <br>
- **Method** - `GET`                  <br>
- **URL Pattern** - `/certficates`  <br>
- **Usage**
    - Open `localhost:8080/certificates` in browser or use Postman
    - **Terminal/CURL**
```
curl -X GET localhost:8080/certificates 
```
- **Expected Response** - All certificates stored
- **NOTE** - There are only 2 certificates in the static data   
- **Example**

![Screenshot](/screenshots/allCertificates.PNG "all certs screenshots")


### 2. Create Certificate
- **Endpoint Name** - `create_certificate`    <br>
- **Method** - `POST`                  <br>
- **URL Pattern** - `/certficates/create`  <br>
- **Usage**
    - Open `localhost:8080/certificates/create` in browser or use Postman
    - **Terminal/CURL**
```
curl -X POST \
  http://localhost:8080/certificates/create \
  -H 'Content-Type: application/json' \
  -H 'OwnerID: rr01' \
  -d '{
        "ID": "c003",
        "Title": "The Yellow House",
        "CreatedAt": "2009-11-17T20:34:58.651387237Z",
        "Year": 1888,
        "Note": "",
        "Transfer": {
            "To": "",
            "Status": ""
        }
}' 
```
- **Expected Response** - Certificate creation successful.
- **NOTE** - Owner of the certificate is expected to be in the header of the request.  The certificate created is also returned on success.
- **Example**

![Screenshot](/screenshots/createCertificate.PNG "status 201: created")



### 3. View Single Certificate
- **Endpoint Name** - `get_certificate`    <br>
- **Method** - `GET`                  <br>
- **URL Pattern** - `/certficates/{id}`  <br>
- **Usage**
    - Open `localhost:8080/certificates/{id}` in browser or use Postman
    - **Terminal/CURL**
```
curl -X GET localhost:8080/certificates/{id}
```
- **Expected Response** - Certificate with the specified ID.
- **NOTE** - Existing IDs in static data: c001, c002
- **Example**

![Screenshot](/screenshots/getCertificate.PNG "status 200: get cert")


### 4. Update Certificate
- **Endpoint Name** - `update_certificate`    <br>
- **Method** - `PUT`                  <br>
- **URL Pattern** - `/certficates/update`  <br>
- **Usage**
    - Open `localhost:8080/certificates/update` in browser or use Postman
    - **Terminal/CURL**
```
curl -X PUT \
  http://localhost:8080/certificates/update \
  -H 'Content-Type: application/json' \
  -d '{
        "ID": "c002",
        "Title": "Café Terrace at Night",
        "CreatedAt": "1888-11-17T20:34:58.651387237Z",
        "OwnerID": "vvg01",
        "Year": 1888,
        "Note": "",
        "Transfer": {
            "To": "",
            "Status": ""
        }
    }'
```
If the certificate does not exist it will be created **IF** the owner is added as a header to the request
- **Expected Response** - Certificate successfully updated.
- **NOTE** - If the new certificate has a different OwnerID than the original, the request will return an appropriate error. <br>*This functionality is done by transfers only.*
- **Example**

![Screenshot](/screenshots/updateCertificate.PNG "status 200: updated")


### 5. Delete Certificate
- **Endpoint Name** - `delete_certificate`    <br>
- **Method** - `DELETE`                  <br>
- **URL Pattern** - `/certficates/delete/{id}`  <br>
- **Usage**
    - Open `localhost:8080/certificates/delete/{id}` in browser or use Postman
    - **Terminal/CURL**
```
curl -X DELETE http://localhost:8080/certificates/delete/c002
```
- **Expected Response** - Certificate successfully deleted.*
- **Example**

![Screenshot](/screenshots/deleteCertificate.PNG "status 200: deleted")


### 6. View List Of User's Certificates
- **Endpoint Name** - `user_certificates`    <br>
- **Method** - `GET`                  <br>
- **URL Pattern** - `/users/{userID}/certificates`  <br>
- **Usage**
    - Open `localhost:8080//users/{userID}/certificates` in browser or use Postman
    - **Terminal/CURL**
```
curl -X GET http://localhost:8080/users/rr01/certificates
```

- **Expected Response** - Return a list of all certificates owned by the user.
- **Example**

![Screenshot](/screenshots/userCertificates.PNG "status 200: user certificates")


### 7. Create Certificate Transfer
- **Endpoint Name** - `create_transfer`    <br>
- **Method** - `POST`                  <br>
- **URL Pattern** - `/certificates/{id}/transfers/create`  <br>
- **Basic Auth Required**
- **Usage**
    - Open `localhost:8080/certificates/{id}/transfers/create` in browser or use Postman
    - **Terminal/CURL**
```
curl -u rr01:rrejh3294 \
-X POST \
  http://localhost:8080/certificates/c001/transfers/create \
  -H 'Content-Type: application/json' \
  -d '{
    "To": "vvg@gmail.com",
    "Status": "pending"
}'
```
Static users data:
```
{
	ID:  "rr01",
	Email:  "reshawnramjattan@gmail.com",
	Name:  "Reshawn",
	password:  "rrejh3294",
},
{
	ID:  "vvg01",
	Email:  "vvg@gmail.com",
	Name:  "Vincent Van Golang",
	password:  "vwh39043f",
}
```
- **Expected Response** - Transfer successfully created.
- **NOTE** - As opposed to just having a custom OwnerID header, transfer operations require basic authorization. <br> Naturally, credentials must belong to the user that owns the certificate.
- **Example**

![Screenshot](/screenshots/transferCreate.PNG "status 201: transfer created")


### 8. Accept Certificate Transfer
- **Endpoint Name** - `accept_transfer`    <br>
- **Method** - `PUT`                  <br>
- **URL Pattern** - `/certificates/{id}/transfers/accept`  <br>
- **Basic Auth Required**
- **Usage**
    - Open `localhost:8080/certificates/{id}/transfers/accept` in browser or use Postman
    - **Terminal/CURL**
```
curl -u vvg01:vwh39043f \
-X PUT http://localhost:8080/certificates/c001/transfers/accept
```
- **Expected Response** - Transfer accepted.
- **NOTE** - Credentials must belong to the user whose email matches the Transfer's To value.
- **Example**

![Screenshot](/screenshots/transferAccept.PNG "status 200: transfer accepted")
