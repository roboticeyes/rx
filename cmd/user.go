package cmd

import (
	"fmt"

	"github.com/breiting/rex"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Get the current user information",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Current user:")
		client := rex.NewClient(ClientID, ClientSecret, nil)
		user := rex.GetCurrentUser(client)
		fmt.Println(user)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
