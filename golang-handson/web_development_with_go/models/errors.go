package models

import "strings"

const (
	ErrNotFound         modelError = "models: resource not found"
	ErrIDInvalid        modelError = "models: ID provided was invalid"
	ErrPasswordInvalid  modelError = "models: incorrect password provided"
	ErrEmailRequired    modelError = "models: email address is required"
	ErrEmailInvalid     modelError = "models: email address is not valid"
	ErrEmailTaken       modelError = "models: email address is already taken"
	ErrPasswordTooShort modelError = "models: password must be at least 8 charaters long"
	ErrPasswordRequired modelError = "models: password is required"
	ErrRememberRequired modelError = "models: remember token is required"
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
	ErrTitleRequired    modelError = "models: title is required"
)

const (
	ErrUserIDRequired privateError = "models: user ID is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
