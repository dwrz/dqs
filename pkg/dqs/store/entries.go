package store

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dwrz/dqs/pkg/dqs/entry"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) DeleteEntry(date string) ( error) {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("dqs bucket does not exist")
		}

		if err := b.Delete([]byte(date)); err != nil {
			return fmt.Errorf(
				"failed to delete bucket entry: %v", err,
			)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetEntry(date string) (*entry.Entry, error) {
	var e *entry.Entry

	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("dqs bucket does not exist")
		}

		data := b.Get([]byte(date))
		if data == nil {
			return nil
		}

		if err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(
			&e,
		); err != nil {
			return fmt.Errorf("failed to gob decode entry: %v", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return e, nil
}

func (s *Store) UpdateEntry(e *entry.Entry) error {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return fmt.Errorf("dqs bucket does not exist")
		}

		var buf bytes.Buffer

		if err := gob.NewEncoder(&buf).Encode(e); err != nil {
			return fmt.Errorf("failed to gob encode entry: %v", err)
		}

		return b.Put([]byte(e.GetKey()), buf.Bytes())
	}); err != nil {
		return err
	}

	return nil
}
