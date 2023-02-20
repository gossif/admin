// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"fmt"

	"github.com/gossif/ebsi"
	"github.com/gossif/admin/wallet"
	"github.com/spf13/cobra"
)

var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Step 3: Register the did document (ebsi only).",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		didString, _ := cmd.Flags().GetString("did")
		did := ebsi.NewDecentralizedIdentifier()
		if err := did.ParseIdentifier(didString); err != nil {
			fmt.Printf("Identifier is not valid\n'%s'\n", err)
			return
		}
		didBucket, err := wallet.GetBucketByDid(did.String())
		if err != nil {
			fmt.Printf("Failed to load the did bucket\n'%s'\n", err)
			return
		}
		didBucket.AdminEncryptionKey, _ = generateSecp256k1AsJwk(didBucket.Did)
		didBucket.AdminTransactionKey, _ = generateSecp256k1AsJwk(didBucket.Did)
		ebsiTrustList := ebsi.NewEBSITrustList(
			ebsi.WithBaseUrl("https://api-pilot.ebsi.eu"),
			ebsi.WithVerbose(true),
		)
		if _, err := ebsiTrustList.RegisterDid(
			ebsi.WithController(didBucket.Did),
			ebsi.WithDocument(didBucket.Document),
			ebsi.WithDocumentMetadata(map[string]interface{}{"deactivated": false}),
			ebsi.WithToken(didBucket.Token),
			ebsi.WithEncryptionKey(didBucket.AdminEncryptionKey),
			ebsi.WithSigningKey(didBucket.AdminSigningKey),
			ebsi.WithTransactionKey(didBucket.AdminTransactionKey),
		); err != nil {
			fmt.Printf("Failed to register the did document\n'%s'\n", err)
			return
		}
		if err = wallet.StoreBucket(didBucket); err != nil {
			fmt.Printf("Failed to save the results\n'%s'\n", err)
			return
		}
		fmt.Printf("Registering of the did document for %s succeeded", didString)
	},
}
