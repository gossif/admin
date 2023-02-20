// Copyright 2023 The Go SSI Framework Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package commands

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// StringPrompt asks for a string value using the label
func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(label)
		fmt.Print(">")

		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

func promptGetAccessToken() string {
	promptToken := stringPrompt("Please provide the access token (https://app-pilot.ebsi.eu/users-onboarding/v2/).")
	re := regexp.MustCompile(`\r?\n`)
	accessToken := re.ReplaceAllString(promptToken, "")
	return accessToken
}
