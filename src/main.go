package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/01100100/forms"
)

const DB_PATH = "/data/db.json"

var REDIRECT_URL = os.Getenv("REDIRECT_URL")

type Form struct {
	name         string
	email        string
	receivedTime time.Time
	FormData     forms.Data
}

func formHandler(response http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/forms" {
		http.Error(response, "404 not found.", http.StatusNotFound)
		return
	}
	if request.Method != "POST" {
		http.Error(response, "Method is not supported.", http.StatusNotFound)
		return
	}

	formData, err := forms.Parse(request)

	fmt.Println("Received request.")

	// Validate
	val := formData.Validator()
	val.Require("name")
	val.Require("email")
	val.MatchEmail("email")
	val.Require("spam")
	val.MatchString("spam", "safe")
	if val.HasErrors() {
		fmt.Println("Received an invalid form")

		// TODO: write back invalid data to client.
		// But do not expose the anti-spam magic ;)
		fmt.Fprintf(response, "Form was invalid!")
		return
	}

	data := &Form{
		name:         formData.Get("name"),
		email:        formData.Get("email"),
		receivedTime: time.Now(),
		FormData:     *formData,
	}

	// Write data to file.
	f, err := os.OpenFile(DB_PATH, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	newLine, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.Write(append(newLine, []byte("\n")...))
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println("Successfully logged form data.")
	http.Redirect(response, request, REDIRECT_URL, http.StatusFound)

	// TODO: process data and send out emails.

	// TODO: serve simple front end.
}

func main() {

	if _, err := os.Stat(DB_PATH); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(DB_PATH)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	http.HandleFunc("/forms", formHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
