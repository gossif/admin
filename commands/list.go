// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"fmt"
	"os"

	"github.com/gossif/admin/wallet"
	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the decentralized identifiers stored in the wallet).",
	Args:  cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		identifiers, err := wallet.GetAllKeys()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while reading the wallet\n'%s'\n", err)
			return
		}
		for _, did := range identifiers {
			fmt.Println(did)
		}
	},
}
