package cmd

import "github.com/spf13/cobra"

var ConfigPath = []string{
	"./configs",
	"../configs",
	"../../configs",
	"../../../configs"}

var WelkomText = `
========================================================================================  
Messaging Services
========================================================================================
- port    : %d
- logrus     : %s
-----------------------------------------------------------------------------------------`

// RootCmd this function for root command
func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "messaging",
		Short: "messaging - messaging Services",
		Long:  "messaging is send message Service",
		Args:  cobra.MinimumNArgs(1),
	}
	return root
}
