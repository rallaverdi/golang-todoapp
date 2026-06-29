package users_postgres_repository

type UserModel struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string // can be null that's why *string
}
