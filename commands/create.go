// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gossif/admin/wallet"
	"github.com/gossif/ebsi"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Step 1: Create an ebsi decentralized identifier.",
	Run: func(_ *cobra.Command, _ []string) {
		did := ebsi.NewDecentralizedIdentifier()
		did.GenerateMethodSpecificId()

		jwkIssuanceKey, _ := generateSecp256r1AsJwk(did.String())
		jwkPublicKey, _ := jwkIssuanceKey.PublicKey()
		didDocument := map[string]interface{}{
			"@context": []string{"https://www.w3.org/ns/did/v1"},
			"id":       did.String(),
			"verificationMethod": []map[string]interface{}{
				{
					"id":           jwkPublicKey.KeyID(),
					"type":         "JsonWebKey2020",
					"controller":   did.String(),
					"publicKeyJwk": jwkPublicKey,
				},
			},
			"authentication":  []string{jwkPublicKey.KeyID()},
			"assertionMethod": []string{jwkPublicKey.KeyID()},
		}
		jwkPresentationKey, _ := generateSecp256r1AsJwk(did.String())
		didBucket := wallet.DidBucket{
			Did:             did.String(),
			Document:        didDocument,
			IssuanceKey:     jwkIssuanceKey,
			PresentationKey: jwkPresentationKey,
		}
		if err := wallet.StoreBucket(didBucket); err != nil {
			fmt.Printf("Failed to save the results\n'%s'\n", err)
			return
		}
		fmt.Printf("Creating of did %s succeeded\n", did.String())
	},
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
