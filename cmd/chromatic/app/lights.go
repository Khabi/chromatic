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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/GetVivid/huego"
)

// lightsCmd represents the lights command
var lightsCmd = &cobra.Command{
	Use:   "lights",
	Short: "List capable lights from hue",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		username, _ := cmd.Flags().GetString("username")

		bridge := huego.New(host, username, "")
		groups, err := bridge.GetEntertainmentGroups()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println("Hue Groups")
		for _, g := range groups {
			fmt.Printf("  %d: %s\n", g.ID, g.Name)
			fmt.Printf("    Light IDs: %s\n", strings.Join(g.Lights, ","))
		}
	},
}

func init() {
	rootCmd.AddCommand(lightsCmd)

	lightsCmd.Flags().StringP("host", "a", "", "Philips Hue hub address")
	lightsCmd.Flags().StringP("username", "u", "", "Philips Hue username")

	lightsCmd.MarkFlagRequired("host")
	lightsCmd.MarkFlagRequired("username")
}
