package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "List available devices",
	Run: func(cmd *cobra.Command, args []string) {
		devices, err := serial.GetPortsList()
		if err != nil {
			log.WithError(err).Error("Error listing devices")
			return
		}

		for _, device := range devices {
			log.WithField("device", device).Info("Found device")
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}
