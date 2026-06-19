package domain

import "errors"

var (
	ErrorTokenExists   = errors.New("Token already exists")
	ErrorUserExists    = errors.New("User already exists")
	ErrorCredentials   = errors.New("Invalid credentials")
	ErrorMonitorExists = errors.New("Monitor already exists")
)
