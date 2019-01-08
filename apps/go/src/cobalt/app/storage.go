package app

import (
	"github.com/boltdb/bolt"
)

const (
	configBucket        = "config"
	tourneyBucket       = "tourney"
	registrationsBucket = "registrations"
)

func makeBucketForPath(tx *bolt.Tx, names ...string) (*bolt.Bucket, error) {
	var err error

	root := []byte(names[0])
	b := tx.Bucket(root)
	if b == nil {
		slog.Debugf("creating bucket for %s", root)
		b, err = tx.CreateBucket(root)
		if err != nil {
			return nil, err
		}
	}
	for _, name := range names[1:] {
		parent := b
		b = parent.Bucket([]byte(name))
		if b == nil {
			slog.Debugf("creating bucket for %s", name)
			b, err = parent.CreateBucket([]byte(name))
			if err != nil {
				return nil, err
			}
		}
	}
	return b, nil
}

func getBucketForPath(tx *bolt.Tx, names ...string) *bolt.Bucket {
	root := []byte(names[0])
	b := tx.Bucket(root)
	for _, name := range names[1:] {
		if b == nil {
			break
		}
		b = b.Bucket([]byte(name))
	}
	return b
}
func getBucketForRegistrations(tx *bolt.Tx, slug string) *bolt.Bucket {
	return getBucketForPath(tx, tourneyBucket, slug, registrationsBucket)
}

func makeBucketForRegistrations(tx *bolt.Tx, slug string) (*bolt.Bucket, error) {
	return makeBucketForPath(tx, tourneyBucket, slug, registrationsBucket)
}

func getBucketForConfig(tx *bolt.Tx) *bolt.Bucket {
	return getBucketForPath(tx, configBucket)
}

func makeBucketForConfig(tx *bolt.Tx) (*bolt.Bucket, error) {
	return makeBucketForPath(tx, configBucket)
}
