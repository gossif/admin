// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"

	"github.com/gossif/admin/commands"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

const (
	Version = "v1.0"
	Banner  = `

    _______     __            __              ___                  _______ __           __ 
   |   _   |.--|  |.--------.|__|.-----.    .'  _|.-----.----.    |    ___|  |--.-----.|__|
   |       ||  _  ||        ||  ||     |    |   _||  _  |   _|    |    ___|  _  |__ --||  |
   |___|___||_____||__|__|__||__||__|__|    |__|  |_____|__|      |_______|_____|_____||__| %s

   Admin for the European Blockchain Services Infrastructure (written in the Go programming language)
   ________________________________________________________________________________________/\__/\__0>___________											
   %s `
)

var rootCmd = &cobra.Command{
	Use:     "essif",
	Version: Version,
	Args:    cobra.ExactArgs(1),
	Short:   "essif - a CLI to manage decentralized identifiers",
}

func init() {
	rootCmd.Long = fmt.Sprintf(Banner, string(colorRed)+Version+string(colorReset), string(colorCyan)+"Implemented by Hietkamp IT-Consultancy"+string(colorReset))

	rootCmd.AddCommand(commands.CreateCmd)
	rootCmd.AddCommand(commands.RegisterCmd)
	rootCmd.AddCommand(commands.OnboardCmd)
	//rootCmd.AddCommand(commands.AccessTokenCmd)
	rootCmd.AddCommand(commands.ResolveCmd)
	rootCmd.AddCommand(commands.ListCmd)

	commands.CreateCmd.Flags().StringP("method", "m", "", "the method used to create the did.")
	commands.CreateCmd.Flags().StringP("domain", "d", "", "the domain for the web method.")
	commands.OnboardCmd.Flags().StringP("did", "d", "", "the did to be onboarded.")
	commands.RegisterCmd.Flags().StringP("did", "d", "", "the did to be registered.")
	//commands.AccessTokenCmd.Flags().StringP("did", "d", "", "the did of the access token")
	commands.ResolveCmd.Flags().StringP("did", "d", "", "the did of the document to resolve")
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("There was an error while executing your CLI\n'%s'\n", err)
		return
	}
}

func main() {
	Execute()
}
