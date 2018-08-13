/*
 * Author: Bernhard Reitinger
 * Date  : 2018
 */

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
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
	paramFile      = "file"
	paramFileList  = "bulk"
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
		client, err := rex.NewClient(RxConfig.ClientID, RxConfig.ClientSecret, nil)
		project, err := rex.GetProjects(client)
		console(err, project)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the details of a project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(RxConfig.ClientID, RxConfig.ClientSecret, nil)

		projectID, _ := cmd.Flags().GetString(paramProjectID)
		project, err := rex.GetProject(client, projectID)

		if c, _ := cmd.Flags().GetInt(paramDownload); c != -1 {
			rex.DownloadFile(client, project.Embedded.ProjectFiles[c].Links.FileDownload.Href)
		} else {
			console(err, project)
		}
	},
}

func getFileEntries(fileName string) []string {

	if fileName == "" {
		return make([]string, 0)
	}
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return make([]string, 0)
	}
	// remove empty entries
	var r []string
	for _, e := range strings.Split(string(content), "\n") {
		if e != "" {
			r = append(r, e)
		}
	}
	return r
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(RxConfig.ClientID, RxConfig.ClientSecret, nil)
		name, _ := cmd.Flags().GetString(paramName)
		fileList, _ := cmd.Flags().GetString(paramFile)

		if files := getFileEntries(fileList); len(files) > 0 {
			for _, f := range files {
				fmt.Print("Creating project ", f, " ... ")
				err = rex.CreateProject(client, f)
				console(err, "Success!")
			}
		} else if name == "" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Project name: ")
			name, _ = reader.ReadString('\n')
			name = strings.Replace(name, "\n", "", -1)
			fmt.Print("Creating new project: ", name)
			err = rex.CreateProject(client, name)
			console(err, "Success!")
		} else {
			fmt.Print("Creating new project: ", name)
			err = rex.CreateProject(client, name)
			console(err, "Success!")
		}
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads a new file for a given project",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := rex.NewClient(RxConfig.ClientID, RxConfig.ClientSecret, nil)
		projectID, _ := cmd.Flags().GetString(paramProjectID)
		localFile, _ := cmd.Flags().GetString(paramLocalFile)
		name := filepath.Base(localFile)
		fileList, _ := cmd.Flags().GetString(paramFileList)

		if cmd.Flag(paramName).Value.String() != "" {
			name, _ = cmd.Flags().GetString(paramName)
		}

		// Bulk upload
		if files := getFileEntries(fileList); len(files) > 0 {
			for _, f := range files {
				fmt.Print("Uploading file ", f, " ... ")
				r, _ := os.Open(f)
				defer r.Close()
				err = rex.UploadProjectFile(client, projectID, filepath.Base(f), f, r)
				console(err, "Success!")
			}
		} else {
			r, _ := os.Open(localFile)
			defer r.Close()
			err = rex.UploadProjectFile(client, projectID, name, localFile, r)
			console(err, "Success!")
		}
	},
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(listCmd)
	projectsCmd.AddCommand(newCmd)
	projectsCmd.AddCommand(uploadFileCmd)
	projectsCmd.AddCommand(showCmd)

	newCmd.Flags().StringP(paramName, "", "", "Name of the project")
	newCmd.Flags().StringP(paramFile, "", "", "File containing a list of project names")

	uploadFileCmd.Flags().StringP(paramProjectID, "", "", "ProjectID [mandatory]")
	uploadFileCmd.Flags().StringP(paramLocalFile, "", "", "Full path of local file [mandatory]")
	uploadFileCmd.Flags().StringP(paramName, "", "", "Name for the project file")
	uploadFileCmd.Flags().StringP(paramFileList, "", "", "File which contains a list of files to be uploaded")

	uploadFileCmd.MarkFlagRequired(paramProjectID)

	showCmd.Flags().StringP(paramProjectID, "", "", "ProjectID [mandatory]")
	showCmd.Flags().Int(paramDownload, -1, "Download the given project file ID")
	showCmd.MarkFlagRequired(paramProjectID)
}
