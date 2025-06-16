package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/timdorr/wattsup/pkg/monitor"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitoring devices",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := setupSignalHandler()

		// Check if we have any devices configured
		if !viper.IsSet("devices") {
			log.Error("No devices configured. Please set up your devices in the configuration file.")
			return
		}

		log.Info("⚡⚡⚡ Starting WattsUp ⚡⚡⚡")

		for name, portFileName := range viper.GetStringMapString("devices") {

			monitor, err := monitor.NewMonitor(name, portFileName)
			if err != nil {
				log.WithField("device", portFileName).WithError(err).Error("Failed to create monitor")
				continue
			}

			go monitor.Start(ctx)
		}

		<-ctx.Done()
		log.Info("⚡⚡⚡ WattsUp stopped ⚡⚡⚡")
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
