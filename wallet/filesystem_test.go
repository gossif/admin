// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package wallet_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/gossif/admin/wallet"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/stretchr/testify/assert"
)

func TestBucket(t *testing.T) {
	t.Run("VerifyNilValues", func(t *testing.T) {
		var (
			expectedDid string = "did:example:123"
		)
		expectedDidBucket := wallet.DidBucket{}
		expectedDidBucket.Did = expectedDid
		expectedDidBucket.IssuanceKey, _ = generateSecp256r1AsJwk(expectedDid)
		expectedDidBucket.PresentationKey, _ = generateSecp256r1AsJwk(expectedDid)

		err := wallet.StoreBucket(expectedDidBucket)
		assert.NoError(t, err)
		_, err = wallet.GetBucketByDid(expectedDid)
		assert.NoError(t, err)
		// change and save the values
	})
	t.Run("StoreBucket", func(t *testing.T) {
		expectedDidBucket := wallet.DidBucket{}
		bucketBytes, _ := os.ReadFile(filepath.Join("testdata", "bucket.json"))
		json.Unmarshal(bucketBytes, &expectedDidBucket)

		err := wallet.StoreBucket(expectedDidBucket)
		assert.NoError(t, err)
	})
	t.Run("GetBucket", func(t *testing.T) {
		var (
			expectedDid string = "did:example:123"
		)
		expectedDidBucket := wallet.DidBucket{}
		bucketBytes, _ := os.ReadFile(filepath.Join("testdata", "bucket.json"))
		err := json.Unmarshal(bucketBytes, &expectedDidBucket)

		assert.NoError(t, err)
		err = wallet.StoreBucket(expectedDidBucket)
		assert.NoError(t, err)

		actualDidBucket, err := wallet.GetBucketByDid(expectedDid)
		assert.NoError(t, err)
		assert.True(t, reflect.DeepEqual(expectedDidBucket.IssuanceKey, actualDidBucket.IssuanceKey))
		assert.True(t, reflect.DeepEqual(expectedDidBucket.PresentationKey, actualDidBucket.PresentationKey))
		assert.True(t, reflect.DeepEqual(expectedDidBucket.AdminEncryptionKey, actualDidBucket.AdminEncryptionKey))
		assert.True(t, reflect.DeepEqual(expectedDidBucket.AdminSigningKey, actualDidBucket.AdminSigningKey))
		assert.True(t, reflect.DeepEqual(expectedDidBucket.Document, actualDidBucket.Document))
		assert.True(t, reflect.DeepEqual(expectedDidBucket.Token, actualDidBucket.Token))
	})

}

// generateSecp256r1AsJwk generates secp256r1 key pair and returns private key as json web key
func generateSecp256r1AsJwk(didController string) (jwk.Key, error) {
	rawKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	kid := didController + "#" + strings.Replace(uuid.NewString(), "-", "", -1)
	jwkKey, err := jwk.FromRaw(rawKey)
	if err != nil {
		return nil, err
	}
	jwkKey.Set(jwk.KeyIDKey, kid)
	return jwkKey, err
}
