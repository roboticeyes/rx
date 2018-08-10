package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/breiting/rex"
	"github.com/spf13/cobra"
)

const (
	paramName      = "name"
	paramProjectID = "id"
	paramLocalFile = "file"
	paramDownload  = "dl"
)

var projectsCmd = &cobra.Command{
	Use:   "project",
	Short: "Work with REX projects",
	Long: `
--------------------------------------------------------------
                          rx - (c) 2018
--------------------------------------------------------------

REX projects are the main entity for organizing your data.

`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List your projects",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(ClientID, ClientSecret, nil)
		project, err := rex.GetProjects(client)
		console(err, project)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the details of a project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(ClientID, ClientSecret, nil)

		projectID, _ := cmd.Flags().GetString(paramProjectID)
		project, err := rex.GetProject(client, projectID)

		if c, _ := cmd.Flags().GetInt(paramDownload); c != -1 {
			rex.DownloadFile(client, project.Embedded.ProjectFiles[c].Links.FileDownload.Href)
		} else {
			console(err, project)
		}
	},
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(ClientID, ClientSecret, nil)
		name, _ := cmd.Flags().GetString(paramName)

		if name == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Project name: ")
			name, _ = reader.ReadString('\n')
			name = strings.Replace(name, "\n", "", -1)
		}

		fmt.Println("Creating new project: ", name)
		err = rex.CreateProject(client, name)
		console(err, "Success!")
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads a new file for a given project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(ClientID, ClientSecret, nil)
		projectID, _ := cmd.Flags().GetString(paramProjectID)
		localFile, _ := cmd.Flags().GetString(paramLocalFile)
		name := filepath.Base(localFile)

		if cmd.Flag(paramName).Value.String() != "" {
			name, _ = cmd.Flags().GetString(paramName)
		}

		r, _ := os.Open(localFile)
		defer r.Close()
		err = rex.UploadProjectFile(client, projectID, name, localFile, r)
		console(err, "Success!")
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listCmd)
	projectsCmd.AddCommand(newCmd)
	projectsCmd.AddCommand(uploadFileCmd)
	projectsCmd.AddCommand(showCmd)

	newCmd.Flags().StringP(paramName, "", "", "Name of the project")

	uploadFileCmd.Flags().StringP(paramProjectID, "", "", "ProjectID [mandatory]")
	uploadFileCmd.Flags().StringP(paramLocalFile, "", "", "Full path of local file [mandatory]")
	uploadFileCmd.Flags().StringP(paramName, "", "", "Name for the project file")

	uploadFileCmd.MarkFlagRequired(paramProjectID)
	uploadFileCmd.MarkFlagRequired(paramLocalFile)

	showCmd.Flags().StringP(paramProjectID, "", "", "ProjectID [mandatory]")
	showCmd.Flags().Int(paramDownload, -1, "Download the given project file ID")
	showCmd.MarkFlagRequired(paramProjectID)

}
