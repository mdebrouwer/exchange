package token

type User struct {
	Name, Email string
}

type Store interface {
	Register(token string, user User) (ok bool, err error)
}
