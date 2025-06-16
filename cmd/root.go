package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "wattsup",
	Short: "A monitoring tool for Sol-Ark 15 inverters",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("wattsup")
	viper.AutomaticEnv()

	viper.SetConfigName("wattsup")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.WithError(err).Error("Error reading config file")
			os.Exit(1)
		}
	}

	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stdout))
}
