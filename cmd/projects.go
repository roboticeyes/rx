package cmd

import (
	"fmt"

	"github.com/breiting/rex"
	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Get the list of projects",
	Long:  `This command gets the list of projects for the current user`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Projects:")
		client := rex.NewClient(ClientID, ClientSecret, nil)
		rex.GetProjects(client)
		// rex.UploadFileContent(client, "/tmp/test.rex") EXPERIMENTA
		// rex.UpdateProjectFile(client) EXPERIMENTAL
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
