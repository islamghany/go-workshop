package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/islamghany/go-workshop/auth/internals/data"
	"github.com/islamghany/go-workshop/auth/internals/validator"
)

func (app *application) hello(w http.ResponseWriter, r *http.Request) {
	// sender := "sender@example.com"
	// subject := "Fancy subject!"
	// body := "Hello from Mailgun Go!"
	// recipient := "ghany1999181@gmail.com"

	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	// fmt.Printf("ID: %s Resp: %s\n", id, resp)

	w.Write([]byte("hello world"))
}

/*
Activating the user process

-As part of the registration process for a new user we will create a cryptographically-secure random activation token that is impossible to guess.
-We will then store a hash of this activation token in a new tokens table, alongside the new user’s ID and an expiry time for the token.
-We will send the original (unhashed) activation token to the user in their welcome email.
-The user subsequently submits their token to a new PUT /v1/users/activated endpoint.
-If the hash of the token exists in the tokens table and hasn’t expired, then we’ll update the activated status for the relevant user to true.
-Lastly, we’ll delete the activation token from our tokens table so that it can’t be used again.
*/
func (app *application) activeUserHandler(w http.ResponseWriter, r *http.Request) {

}
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the data from the request body into a new User struct. Notice also that we
	// set the Activated field to false, which isn't strictly necessary because the
	// Activated field will have the zero-value of false by default. But setting this
	// explicitly helps to make our intentions clear to anyone reading the code.
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.background(func() {
		sender := "auth@example.com"
		subject := "Activate Your Account!"
		body := `
		<h4>your acount activated link is: <a href="http://localhost/8000/users/activate/%s">here</a></h4>
		<p>http://localhost/8000/users/activate/%s</p>
		`

		_, _, err := app.sendEmail(sender, subject, fmt.Sprintf(body, token.Plaintext, token.Plaintext), user.Email)

		if err != nil {
			log.Println(err)
		}

	})

	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
