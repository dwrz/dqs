package store

import (
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	// One bucket stores everything.
	// The "user" key maps to the user data.
	// YYYYMMDD date keys map to entries.
	bucketName        = "dqs"
	dbFile            = "dqs.db"
	dbDirPemissions   = 0700
	dbFilePermissions = 0600
)

type Store struct {
	db *bolt.DB
}

func Open(path string) (*Store, error) {
	var store Store

	// Create the path to the DB, if it doesn't exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, dbDirPemissions); err != nil {
			return nil, fmt.Errorf(
				"failed to create dqs directory: %v", err,
			)
		}
	}

	db, err := bolt.Open(path+dbFile, dbFilePermissions, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %v", err)
	}

	store.db = db

	if err := store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to create bucket: %v", err)
	}

	return &store, nil
}

func (s *Store) Close() {
	s.db.Close()
}
