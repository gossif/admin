// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"fmt"
	"strings"

	secp256k1v4 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/google/uuid"
	"github.com/gossif/admin/wallet"
	"github.com/gossif/ebsi"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/spf13/cobra"
)

var OnboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Step 2: Onboard the controller of the did.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		didString, _ := cmd.Flags().GetString("did")
		did := ebsi.NewDecentralizedIdentifier()
		if err := did.ParseIdentifier(didString); err != nil {
			fmt.Printf("Identifier is not valid\n'%s'\n", err)
			return
		}
		accessToken := promptGetAccessToken()
		didBucket, err := wallet.GetBucketByDid(did.String())
		if err != nil {
			fmt.Printf("Failed to load the did bucket\n'%s'\n", err)
			return
		}
		didBucket.AdminSigningKey, _ = generateSecp256k1AsJwk(didBucket.Did)
		ebsiTrustList := ebsi.NewEBSITrustList(
			ebsi.WithBaseUrl("https://api-pilot.ebsi.eu"),
			ebsi.WithVerbose(true),
			ebsi.WithAuthToken(accessToken),
		)
		// token is a capthca token or a vc jwt
		token, err := ebsiTrustList.Onboard(did.String(), didBucket.AdminSigningKey)
		if err != nil {
			fmt.Printf("failed to onboard the user\n'%s'\n", err)
			return
		}
		switch token := token.(type) {
		case string:
			didBucket.Token = token
			if err = wallet.StoreBucket(didBucket); err != nil {
				fmt.Printf("failed to save the results\n'%s'\n", err)
				return
			}
		default:
			fmt.Printf("invalid response type\n'%s'\n", err)
		}

		fmt.Printf("Onboarding of %s succeeded\n", didString)
	},
}

// generateSecp256k1AsJwk generates secp256k1 key pair and returns private key as json web key
func generateSecp256k1AsJwk(didController string) (jwk.Key, error) {
	rawKey, err := secp256k1v4.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	kid := didController + "#" + strings.Replace(uuid.NewString(), "-", "", -1)
	jwkKey, err := jwk.FromRaw(rawKey.ToECDSA())
	if err != nil {
		return nil, err
	}
	jwkKey.Set(jwk.KeyIDKey, kid)
	return jwkKey, err
}
