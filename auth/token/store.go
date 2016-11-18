package token

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Store interface {
	Register(token []byte, user User) (ok bool, err error)
	Get(token []byte) (user User, ok bool, err error)
}
