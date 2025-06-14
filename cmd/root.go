package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"

	"github.com/joho/godotenv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "wattsup",
	Short: "A monitoring tool for Sol-Ark 15 inverters",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := setupSignalHandler()

		log.Info("⚡⚡⚡ Starting WattsUp ⚡⚡⚡")

		<-ctx.Done()
		log.Info("⚡⚡⚡ WattsUp stopped ⚡⚡⚡")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if _, err := os.Stat(".env"); err == nil {
		log.Info("Loading environment file")
		err := godotenv.Load()
		if err != nil {
			log.WithError(err).Fatal("Error loading environment file")
		}
	}

	viper.SetEnvPrefix("wattsup")
	viper.AutomaticEnv()

	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stdout))
}
