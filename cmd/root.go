package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/rai-project/cmd"
	"github.com/rai-project/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	isColor   bool
	isVerbose bool
	isDebug   bool
)

var RootCmd = &cobra.Command{
	Use: "inle",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initColor)

	RootCmd.AddCommand(cmd.VersionCmd)
	RootCmd.AddCommand(cmd.LicenseCmd)
	RootCmd.AddCommand(cmd.EnvCmd)
	RootCmd.AddCommand(cmd.GendocCmd)
	RootCmd.AddCommand(cmd.CompletionCmd)
	RootCmd.AddCommand(cmd.BuildTimeCmd)

	RootCmd.PersistentFlags().BoolVarP(&isColor, "color", "c", true, "Toggle color output.")
	RootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Toggle verbose mode.")
	RootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "d", false, "Toggle debug mode.")

	viper.BindPFlag("app.debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("app.verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("app.color", RootCmd.PersistentFlags().Lookup("color"))
}

func initConfig() {
	opts := []config.Option{
		config.AppName("inle"),
		config.ConfigString(""),
		config.DebugMode(isDebug),
		config.DebugMode(isVerbose),
		config.ColorMode(isColor),
	}
	config.Init(opts...)
}

func initColor() {
	color.NoColor = !isColor
}
