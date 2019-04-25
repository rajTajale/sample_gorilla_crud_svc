package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var accounts []Account

// Account is used to hold the account details
type Account struct {
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	MobileNumber int64  `json:"mobilenumber"`
	Password     string `json:"password"`
}

// CreateAccount is used to create an account
func CreateAccount(w http.ResponseWriter, r *http.Request) {

	account := &Account{}

	// read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshal body into the object called account
	err = json.Unmarshal(body, account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate account details
	for _, acc := range accounts {
		if acc.FirstName == account.FirstName {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Account already exist"))
			return
		}
	}

	// store account details in an account
	accounts = append(accounts, *account)

	msg := fmt.Sprintf("Hello %s, welcome to our new application", account.FirstName)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(msg))
}

// ListAccounts is used to list all the accounts
func ListAccounts(w http.ResponseWriter, r *http.Request) {

	resp, err := json.Marshal(accounts)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetDetailsByName is used to get the details of an account by name
func GetDetailsByName(w http.ResponseWriter, r *http.Request) {

	queryParams := mux.Vars(r)

	for _, acc := range accounts {

		if acc.FirstName == queryParams["name"] {

			resp, err := json.Marshal(acc)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(resp)
			return
		}

	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Account not found"))

}

// DeleteAccount is used to detete the details of an account by name
func DeleteAccount(w http.ResponseWriter, r *http.Request) {

	var isDeleted bool

	a := []Account{}

	queryParams := mux.Vars(r)

	for _, acc := range accounts {

		if acc.FirstName != queryParams["name"] {
			a = append(a, acc)
		} else {
			isDeleted = true
		}
	}

	accounts = a

	if !isDeleted {
		w.Write([]byte("Account not found"))
	} else {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("Account deleted"))
	}

}

// UpdateAccount is used to update the account details
func UpdateAccount(w http.ResponseWriter, r *http.Request) {

	var isUpdated bool

	a := []Account{}
	account := &Account{}

	queryParams := mux.Vars(r)

	for _, acc := range accounts {

		if acc.FirstName == queryParams["name"] {

			// read request body
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)
			}

			// unmarshal body into the object called account
			err = json.Unmarshal(body, account)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			acc.FirstName = account.FirstName
			acc.LastName = account.LastName
			acc.MobileNumber = account.MobileNumber
			acc.Password = account.Password

			isUpdated = true

		}

		a = append(a, acc)
	}

	if isUpdated {
		accounts = a

		msg := fmt.Sprintf("Hello %v, your account is updated", account.FirstName)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	} else {
		msg := fmt.Sprintln("Account doesnot exist")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
	}

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/account", CreateAccount).Methods(http.MethodPost)          // to create an account
	router.HandleFunc("/account", ListAccounts).Methods(http.MethodGet)            // to get the details of all accounts
	router.HandleFunc("/account/{name}", GetDetailsByName).Methods(http.MethodGet) // to get the details of specific account
	router.HandleFunc("/account/{name}", DeleteAccount).Methods(http.MethodDelete) // to delete the details of specific account
	router.HandleFunc("/account/{name}", UpdateAccount).Methods(http.MethodPut)    //to update the account details

	fmt.Println("listening started in : 8080")

	http.ListenAndServe(":8080", router)
}
