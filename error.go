package main

import "fmt"

type NotAuthenticatedError string

func (e NotAuthenticatedError) Error() string {
	return fmt.Sprintf("Can't make request; the user is not authenticated:\n%v\n", string(e))
}

type BadRequestError string

func (e BadRequestError) Error() string {
	return fmt.Sprintf("Request Failed due to incorrect formatting:\n%v\n", string(e))
}

type EmailExistsError string

func (e EmailExistsError) Error() string {
	return fmt.Sprintf("Request failed due to the provided email already being registered in the database:\n%v\n", string(e))
}

type InvalidCredentialsError string

func (e InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Request failed due to the provided credentials being invalid/incorrect:\n%v\n", string(e))
}
