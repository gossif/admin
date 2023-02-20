// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"encoding/json"
	"fmt"

	"github.com/gossif/ebsi"
	"github.com/spf13/cobra"
)

var ResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Step 4: Resolve a did document.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, _ []string) {
		didString, _ := cmd.Flags().GetString("did")
		did := ebsi.NewDecentralizedIdentifier()
		if err := did.ParseIdentifier(didString); err != nil {
			fmt.Printf("Identifier is not valid\n'%s'\n", err)
			return
		}
		ebsiTrustList := ebsi.NewEBSITrustList(
			ebsi.WithBaseUrl("https://api-pilot.ebsi.eu"),
		)
		rawdoc, err := ebsiTrustList.ResolveDid(did.String())
		if err != nil {
			fmt.Printf("Failed to resolve the did document.\n%s\n", err)
			return
		}
		jsonDiddoc, _ := json.MarshalIndent(rawdoc, "", "    ")
		fmt.Printf("Resolving the did document succeeded.\n%s\n", string(jsonDiddoc))
	},
}
