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
	Long: getCmdLineHeader() + `
REX user is the identity management console.

`,
	Run: func(cmd *cobra.Command, args []string) {

		if c, _ := cmd.Flags().GetBool(paramCurrentUser); c == true {

			console(nil, RxConfig.Client.User)

		} else if c, _ := cmd.Flags().GetBool(paramCount); c == true {
			count, err := rex.GetTotalNumberOfUsers(RxConfig.Client)
			console(err, fmt.Sprintf("Found %d registered users.\n", count))

		} else if email, _ := cmd.Flags().GetString(paramFindUser); email != "" {

			user, err := rex.GetUserByEmail(RxConfig.Client, email)
			console(err, user)

		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	usersCmd.Flags().BoolP(paramCurrentUser, "m", false, "Show your user information")
	usersCmd.Flags().BoolP(paramCount, "c", false, "Total amount of users")

	usersCmd.Flags().StringP(paramFindUser, "", "", "Find a user by email")
}
