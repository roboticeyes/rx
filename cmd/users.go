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
	Short: "Get user(s) information",
	Run: func(cmd *cobra.Command, args []string) {

		client := rex.NewClient(ClientID, ClientSecret, nil)

		if c, _ := cmd.Flags().GetBool(paramCurrentUser); c == true {

			user, err := rex.GetCurrentUser(client)
			console(err, user)

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

func console(err error, value interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\n%v\n", value)
	}
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.Flags().BoolP(paramCurrentUser, "m", false, "Show your user information")
	usersCmd.Flags().BoolP(paramCount, "c", false, "Total amount of users")

	usersCmd.Flags().StringP(paramFindUser, "", "", "Find a user by email")
}
