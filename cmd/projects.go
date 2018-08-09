package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/breiting/rex"
	"github.com/spf13/cobra"
)

const (
	paramList   = "list"
	paramCreate = "create"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Work with REX projects",
	Long:  `Create and work with your personal REX projects`,
	Run: func(cmd *cobra.Command, args []string) {

		client, err := rex.NewClient(ClientID, ClientSecret, nil)
		if err != nil {
			panic(err)
		}

		if c, _ := cmd.Flags().GetBool(paramList); c == true {
			project, err := rex.GetProjects(client)
			console(err, project)
		} else if c, _ := cmd.Flags().GetBool(paramCreate); c == true {

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Project name: ")
			name, _ := reader.ReadString('\n')
			name = strings.Replace(name, "\n", "", -1)
			fmt.Println("Creating new project: ", name)
			err := rex.CreateProject(client, name)
			console(err, "Success!")
		}

		// rex.UploadFileContent(client, "/tmp/test.rex") EXPERIMENTA
		// rex.UpdateProjectFile(client) EXPERIMENTAL
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)

	projectsCmd.Flags().BoolP(paramList, "l", false, "List your projects [default]")
	projectsCmd.Flags().BoolP(paramCreate, "c", false, "Create a new project")
}
