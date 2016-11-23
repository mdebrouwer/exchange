package token

import (
	"encoding/json"

	"github.com/boltdb/bolt"
)

type Token []byte

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Store interface {
	Register(token Token, user User) (ok bool, err error)
	Get(token Token) (user User, ok bool, err error)
}

type boltBackedStore struct {
	db *bolt.DB
}

func NewBoltBackedStore(path string) (Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &boltBackedStore{db: db}, nil
}

func (s *boltBackedStore) Get(token Token) (user User, ok bool, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("USERS"))
		if bucket == nil {
			return nil
		}
		data := bucket.Get(token)
		ok = (data != nil)
		if ok {
			return json.Unmarshal(data, &user)
		}
		return nil
	})
	return
}

func (s *boltBackedStore) Register(token Token, user User) (ok bool, err error) {
	err = s.db.Update(func(tx *bolt.Tx) error {
		bucket, e := tx.CreateBucketIfNotExists([]byte("USERS"))
		if e != nil {
			return e
		}
		ok = (bucket.Get(token) == nil)
		if !ok {
			return nil
		}
		data, e := json.Marshal(user)
		if e != nil {
			return e
		}
		return bucket.Put(token, data)
	})
	return
}
