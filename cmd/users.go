package cmd

import (
	"fmt"

	"github.com/breiting/rex"
	"github.com/spf13/cobra"
)

const (
	paramCurrentUser = "me"
	paramCount       = "count"
	paramFindUser    = "findByEmail"
)

// usersmd represents the user command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Work with REX users",
	Run: func(cmd *cobra.Command, args []string) {

		client, err := rex.NewClient(ClientID, ClientSecret, nil)
		if err != nil {
			panic(err)
		}

		if c, _ := cmd.Flags().GetBool(paramCurrentUser); c == true {

			console(err, client.User)

		} else if c, _ := cmd.Flags().GetBool(paramCount); c == true {
			count, err := rex.GetTotalNumberOfUsers(client)
			console(err, fmt.Sprintf("Found %d registered users.\n", count))

		} else if email, _ := cmd.Flags().GetString(paramFindUser); email != "" {

			user, err := rex.GetUserByEmail(client, email)
			console(err, user)

		} else {
			fmt.Println("Flag required")
		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.Flags().BoolP(paramCurrentUser, "m", false, "Show your user information")
	usersCmd.Flags().BoolP(paramCount, "c", false, "Total amount of users")

	usersCmd.Flags().StringP(paramFindUser, "", "", "Find a user by email")
}
