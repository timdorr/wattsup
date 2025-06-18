package cmd

import (
	"os"
	"time"

	"github.com/timdorr/wattsup/pkg/config"
	"github.com/timdorr/wattsup/pkg/monitor"

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
	Run: func(cmd *cobra.Command, args []string) {
		ctx := setupSignalHandler()

		config := config.GetConfig()

		log.Info("⚡⚡⚡ Starting WattsUp ⚡⚡⚡")

		for _, device := range config.Devices {
			log.WithField("device", device).Infof("Starting monitor for device: %s", device.Name)

			monitor := monitor.NewMonitor(device.Name, device.File, device.ID, config.Registers, config.Database)
			go monitor.Start(ctx)
		}

		<-ctx.Done()
		time.Sleep(500 * time.Millisecond)
		log.Info("⚡⚡⚡ WattsUp stopped ⚡⚡⚡")
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
	log.SetLevel(log.DebugLevel)
	log.SetHandler(text.New(os.Stdout))
}
