/*
Copyright © 2020 Richard Cox <code@bot37.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/GetVivid/huego"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Create a user on the philips hue bridge",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")

		bridge := huego.Bridge{
			Host: host,
		}

		fmt.Println("Press the link button on the hub, then press enter...")
		fmt.Scanln()

		user, clientkey, err := bridge.CreateUser("vivid ambient lights") // Link button needs to be pressed
		if err != nil {
			var e *huego.APIError
			if errors.As(err, &e) {
				fmt.Fprintf(os.Stderr, "Error creating user: %s\n", e.Description)
			} else {
				fmt.Fprintf(os.Stderr, "Error creating user: %s\n", err.Error())
			}
			os.Exit(1)
		}

		fmt.Printf("Username: %s\n", user)
		fmt.Printf("ClientKey: %s\n", clientkey)

	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	registerCmd.Flags().StringP("host", "a", "", "Philips Hue hub address")
	registerCmd.MarkFlagRequired("host")
}
