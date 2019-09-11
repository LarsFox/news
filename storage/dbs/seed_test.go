package dbs

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

// TestSeed создает БД и помещает в нее две тестовые статьи.
func TestSeed(t *testing.T) {
	db, err := bolt.Open("../.tmp/seed.db", 0600, nil)
	if err != nil {
		t.Error(err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketNews)
		if err != nil {
			return err
		}

		if err := putInBucket(
			"how-to-seed",
			&NewsPiece{Header: "How to seed", Date: time.Now().Unix()},
			bucket,
		); err != nil {
			return err
		}
		if err := putInBucket(
			"how-not-to-seed",
			&NewsPiece{Header: "How not to seed", Date: time.Now().Unix()},
			bucket,
		); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
}

func putInBucket(key string, news *NewsPiece, bucket *bolt.Bucket) error {
	b, err := json.Marshal(news)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), b)
}
