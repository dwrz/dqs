package store

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dwrz/dqs/pkg/dqs/user"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) GetUser() (*user.User, error) {
	var u *user.User

	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("dqs bucket does not exist")
		}

		data := b.Get([]byte("user"))
		if data == nil {
			// If a user hasn't been set, return the default user.
			defaultUser := user.DefaultUser
			u = &defaultUser
			return nil
		}

		if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(
			&u,
		); err != nil {
			return fmt.Errorf("failed to gob decode user: %v", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to get user: %v\n", err)
	}

	return u, nil
}

func (s *Store) UpdateUser(u *user.User) error {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("dqs bucket does not exist")
		}

		var buf bytes.Buffer

		if err := gob.NewEncoder(&buf).Encode(u); err != nil {
			return fmt.Errorf("failed to gob encode user: %v", err)
		}

		return b.Put([]byte("user"), buf.Bytes())
	}); err != nil {
		return err
	}

	return nil
}
