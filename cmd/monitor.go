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

		for name, device := range viper.GetStringMap("devices") {
			log.WithField("device", device).Infof("Starting monitor for device: %s", name)

			portFileName := device.(map[string]interface{})["device"].(string)
			slaveId := int(device.(map[string]interface{})["slaveid"].(float64))

			monitor := monitor.NewMonitor(name, portFileName, slaveId)
			go monitor.Start(ctx)
		}

		<-ctx.Done()
		log.Info("⚡⚡⚡ WattsUp stopped ⚡⚡⚡")
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
