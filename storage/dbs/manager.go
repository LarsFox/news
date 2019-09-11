package dbs

import (
	"encoding/json"
	"log"

	"github.com/boltdb/bolt"
)

// Manager ...
type Manager struct {
	db *bolt.DB
}

// NewManager - конструктор.
// Создает БД, если ее нет, и проводит необходимые операции в ней.
func NewManager(dbPath string) (*Manager, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketNews); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &Manager{db: db}, nil
}

func (m *Manager) GetNewsPiece(newsID string) (*NewsPiece, error) {
	newsPiece := &NewsPiece{}
	if err := m.db.View(func(tx *bolt.Tx) error {
		val := tx.Bucket(bucketNews).Get([]byte(newsID))
		if val == nil {
			return ErrNotFound
		}
		if err := json.Unmarshal(val, newsPiece); err != nil {
			log.Println(err)
			return ErrInternal
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return newsPiece, nil
}
