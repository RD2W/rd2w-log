package repository

// UserRepository combines all user operations
type UserRepository interface {
	UserReader
	UserWriter
}
