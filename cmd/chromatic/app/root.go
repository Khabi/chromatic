/*
Copyright © 2020 Richard Cox

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
	"path"
	"regexp"
	"strconv"

	"github.com/GetVivid/huego"
	"github.com/Khabi/chromatic/internal/api"
	"github.com/Khabi/chromatic/internal/chromatic"
	"github.com/Khabi/chromatic/internal/location"
	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chromatic",
	Short: "Ambient lighting controller",
	Long:  `A personal ambient lighting controller.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		lvl, err := logrus.ParseLevel(viper.GetString("log_level"))
		if err != nil {
			logrus.Warn("invalid log level, setting to info")
			lvl = logrus.InfoLevel
		}
		logrus.SetLevel(lvl)
		commandChan := make(chan chromatic.State)
		statusChan := make(chan chromatic.ServerStatus)

		// Configure the video device
		if viper.GetString("video.device") == "" || viper.GetString("video.profile") == "" {
			fmt.Fprintln(os.Stderr, "video input misconfigured")
			os.Exit(255)
		}

		video, err := v4l.Open(viper.GetString("video.device"))
		if err != nil {
			fmt.Println("Unable to open video device.")
			os.Exit(1)
		}

		cfg, err := video.GetConfig()
		if err != nil {
			fmt.Println("Video profile issues")
			os.Exit(1)
		}

		cfg.Format = mjpeg.FourCC

		re := regexp.MustCompile(`(?P<width>\d+)x(?P<height>\d+)@(?P<fps>\d+)`)
		profile := re.FindStringSubmatch(viper.GetString("video.profile"))

		width, err := strconv.Atoi(profile[1])
		if err != nil {
			fmt.Println("invalid profile width")
			os.Exit(1)
		}
		height, err := strconv.Atoi(profile[2])
		if err != nil {
			fmt.Println("invalid profile height")
			os.Exit(1)
		}
		fps, err := strconv.Atoi(profile[3])
		if err != nil {
			fmt.Println("invalid profile fps")
			os.Exit(1)
		}
		cfg.Width = width
		cfg.Height = height
		cfg.FPS = v4l.Frac{uint32(fps), 1}
		err = video.SetConfig(cfg)
		if err != nil {
			fmt.Println(err)
			fmt.Println("invalid video configuration")
			os.Exit(1)
		}

		//Configure Hue
		bridge := huego.New(
			viper.GetString("light.bridge"),
			viper.GetString("light.username"),
			viper.GetString("light.client_key"),
		)
		var group *huego.EntertainmentGroup
		if viper.GetInt("light.group_id") != 0 {
			group, err = bridge.GetEntertainmentGroup(viper.GetInt("light.group_id"))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		if viper.GetString("light.group_name") != "" {
			groups, err := bridge.GetEntertainmentGroups()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, g := range groups {
				if g.Name == viper.GetString("light.group_name") {
					group = &g
				}
			}
		}

		if group == nil {
			fmt.Println("no matching entertainment group")
			os.Exit(1)
		}

		// Get bounds for light sources
		var bounds location.Bounds
		for id, loc := range group.Locations {
			preset := viper.GetString(fmt.Sprintf("light.binding.%d", id))
			switch preset {
			case "top":
				bounds = append(bounds, location.Preset(id, location.Top))
			case "left":
				bounds = append(bounds, location.Preset(id, location.Left))
			case "bottom":
				bounds = append(bounds, location.Preset(id, location.Bottom))
			case "right":
				bounds = append(bounds, location.Preset(id, location.Right))
			case "whole":
				bounds = append(bounds, location.Preset(id, location.Whole))
			default:
				bounds = append(bounds, location.Bound{ID: id, X: loc.X, Y: loc.Y, Width: 20, Height: 20})
			}
		}

		go chromatic.Run(commandChan, statusChan, video, group, bounds)

		api.Run(viper.GetString("bind"), commandChan, statusChan)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/chromatic.yaml or ~/.config)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".chromatic" (without extension).
		viper.AddConfigPath(path.Join(home, ".config"))
		viper.AddConfigPath("/etc/")
		viper.SetConfigName("chromatic")
	}

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Missing config file!")
		os.Exit(1)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())
}
