package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/01100100/forms"
)

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
	f, err := os.OpenFile("db.json", os.O_APPEND|os.O_WRONLY, 0644)
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
	fmt.Fprintf(response, "form sent successfully!")
	fmt.Println("Successfully logged form data.")

	// TODO: process data and send out emails.

	// TODO: serve simple front end.
}

func main() {
	http.HandleFunc("/forms", formHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
