// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package wallet

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/tidwall/buntdb"
)

type DidBucket struct {
	Did                 string                 `json:"did,omitempty"`
	IssuanceKey         jwk.Key                `json:"issKey,omitempty"`
	PresentationKey     jwk.Key                `json:"presKey,omitempty"`
	AdminEncryptionKey  jwk.Key                `json:"encKey,omitempty"`
	AdminTransactionKey jwk.Key                `json:"txnKey,omitempty"`
	AdminSigningKey     jwk.Key                `json:"sigKey,omitempty"`
	Document            map[string]interface{} `json:"doc,omitempty"`
	Token               string                 `json:"token,omitempty"`
}

type rawDidBucket struct {
	Did                 string          `json:"did,omitempty"`
	IssuanceKey         json.RawMessage `json:"issKey,omitempty"`
	PresentationKey     json.RawMessage `json:"presKey,omitempty"`
	AdminEncryptionKey  json.RawMessage `json:"encKey,omitempty"`
	AdminTransactionKey json.RawMessage `json:"txnKey,omitempty"`
	AdminSigningKey     json.RawMessage `json:"sigKey,omitempty"`
	Document            json.RawMessage `json:"doc,omitempty"`
	Token               string          `json:"token,omitempty"`
}

var dbStore *MKVStore

func init() {
	var (
		err error
	)
	dbStore, err = NewFileKVStore()
	if err != nil {
		panic(err)
	}
}

// MemoryStore token storage based on buntdb(https://github.com/tidwall/buntdb)
type MKVStore struct {
	*buntdb.DB
}

// NewMemoryKVStore create a store instance based on memory
func NewMemoryKVStore() (*MKVStore, error) {
	db, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, err
	}
	return &MKVStore{db}, nil
}

// NewMemoryKVStore create a store instance based on a file
func NewFileKVStore() (*MKVStore, error) {
	db, err := buntdb.Open("walletdata.db")
	if err != nil {
		return nil, err
	}
	return &MKVStore{db}, nil
}

// Set persist the value with key
func (m *MKVStore) Set(key string, value string, expires time.Duration) error {
	var (
		expiresOption *buntdb.SetOptions = nil
	)
	if strings.TrimSpace(key) == "" {
		return errors.New("key is empty")
	}
	m.DB.Update(func(tx *buntdb.Tx) error {
		if expires > 0 {
			// add 5 seconds for processing time
			expires += time.Second * 5
			expiresOption = &buntdb.SetOptions{Expires: true, TTL: expires}
		}
		tx.Set(key, value, expiresOption)
		return nil
	})
	return nil
}

func (m *MKVStore) Get(key string) (string, error) {
	var (
		value string
	)
	if strings.TrimSpace(key) == "" {
		return "", buntdb.ErrNotFound
	}
	err := m.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		value = val
		return nil
	})
	if err != nil {
		return "", err
	}
	return value, nil
}

// remove key
func (m *MKVStore) Remove(key string) error {
	if strings.TrimSpace(key) == "" {
		return buntdb.ErrNotFound
	}
	err := m.DB.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
	return err
}

func (m *MKVStore) GetAllKeys() ([]string, error) {
	var allKeys []string
	err := m.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, _ string) bool {
			allKeys = append(allKeys, key)
			return true // continue iteration
		})
		return err
	})
	if err != nil {
		return []string{}, nil
	}
	return allKeys, nil
}

func StoreBucket(bucket DidBucket) error {
	bucketBytes, err := bucket.MarshalJSON()
	if err != nil {
		return err
	}
	return dbStore.Set(bucket.Did, string(bucketBytes), -1)
}

func GetBucketByDid(didSubject string) (DidBucket, error) {
	bucketString, err := dbStore.Get(didSubject)
	if err != nil {
		return DidBucket{}, err
	}
	bucket := DidBucket{}
	if err = json.Unmarshal([]byte(bucketString), &bucket); err != nil {
		return DidBucket{}, err
	}
	if bucket.Did != didSubject {
		return bucket, errors.New("not found")
	}
	return bucket, nil
}

func GetAllKeys() ([]string, error) {
	return dbStore.GetAllKeys()
}

func (bucket *DidBucket) MarshalJSON() ([]byte, error) {
	rawBucket := rawDidBucket{}

	elements := reflect.ValueOf(bucket).Elem()
	for i := 0; i < elements.NumField(); i++ {
		switch element := elements.Field(i).Interface().(type) {
		case string:
			switch elements.Type().Field(i).Name {
			case "Did":
				rawBucket.Did = element
			case "Token":
				rawBucket.Token = element
			}
		case jwk.Key:
			if element != nil {
				switch elements.Type().Field(i).Name {
				case "IssuanceKey":
					rawBucket.IssuanceKey, _ = json.Marshal(element)
				case "PresentationKey":
					rawBucket.PresentationKey, _ = json.Marshal(element)
				case "AdminEncryptionKey":
					rawBucket.AdminEncryptionKey, _ = json.Marshal(element)
				case "AdminTransactionKey":
					rawBucket.AdminTransactionKey, _ = json.Marshal(element)
				case "AdminSigningKey":
					rawBucket.AdminSigningKey, _ = json.Marshal(element)
				}
			}
		case map[string]interface{}:
			if element != nil {
				switch elements.Type().Field(i).Name {
				case "Document":
					rawBucket.Document, _ = json.Marshal(element)
				}
			}
		}
	}
	return json.Marshal(rawBucket)
}

func (bucket *DidBucket) UnmarshalJSON(data []byte) error {
	var (
		err error
	)
	rawBucket := rawDidBucket{}
	if err = json.Unmarshal(data, &rawBucket); err != nil {
		return err
	}
	elements := reflect.ValueOf(&rawBucket).Elem()
	for i := 0; i < elements.NumField(); i++ {
		switch element := elements.Field(i).Interface().(type) {
		case string:
			switch elements.Type().Field(i).Name {
			case "Did":
				bucket.Did = element
			case "Token":
				bucket.Token = element
			}
		case json.RawMessage:
			if element != nil {
				switch elements.Type().Field(i).Name {
				case "IssuanceKey":
					bucket.IssuanceKey, err = jwk.ParseKey(element)
				case "PresentationKey":
					bucket.PresentationKey, err = jwk.ParseKey(element)
				case "AdminEncryptionKey":
					bucket.AdminEncryptionKey, err = jwk.ParseKey(element)
				case "AdminTransactionKey":
					bucket.AdminTransactionKey, err = jwk.ParseKey(element)
				case "AdminSigningKey":
					bucket.AdminSigningKey, err = jwk.ParseKey(element)
				case "Document":
					err = json.Unmarshal(element, &bucket.Document)
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
