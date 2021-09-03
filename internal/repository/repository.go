package repository

type DatabaseRepo interface {

	// Implemented in postgres.go file
	AllUsers() bool
}