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
	"strconv"
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

	openStreetMap       = "https://www.openstreetmap.org"
	openStreetMapSearch = "https://nominatim.openstreetmap.org/search"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
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
		project, err := rex.GetProjects(RxConfig.Client, RxConfig.Client.User.UserID)
		console(err, project)
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the details of a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, _ := cmd.Flags().GetString(paramProjectID)
		project, err := rex.GetProject(RxConfig.Client, projectID)

		if c, _ := cmd.Flags().GetInt(paramDownload); c != -1 {
			rex.DownloadFile(RxConfig.Client, project.Embedded.ProjectFiles[c].Links.FileDownload.Href)
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
		name, _ := cmd.Flags().GetString(paramName)
		fileList, _ := cmd.Flags().GetString(paramFile)

		if files := getFileEntries(fileList); len(files) > 0 {
			for _, f := range files {
				fmt.Print("Creating project ", f, " ... ")
				err := rex.CreateProject(RxConfig.Client, RxConfig.Client.User.UserID, f, nil, nil)
				console(err, "Success!")
			}
		} else if name == "" {
			createProjectInteractive()
		} else {
			fmt.Print("Creating new project: ", name)
			err := rex.CreateProject(RxConfig.Client, RxConfig.Client.User.UserID, name, nil, nil)
			console(err, "Success!")
		}
	},
}

var uploadFileCmd = &cobra.Command{
	Use:   "upload",
	Short: "Uploads a new file for a given project",
	Run: func(cmd *cobra.Command, args []string) {
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
				err := rex.UploadProjectFile(RxConfig.Client, projectID, filepath.Base(f), f, nil, r)
				console(err, "Success!")
			}
		} else {
			uploadProjectFileInteractive(projectID, name, localFile)
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

func getInput(name string, r *bufio.Reader) string {

	fmt.Printf("%-20s: ", name)
	name, _ = r.ReadString('\n')
	return strings.Replace(name, "\n", "", -1)
}

func interpretFloatInput(input string, val float64) float64 {
	if input != "" {
		l, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0.0
		}
		return l
	}
	return val
}

func uploadProjectFileInteractive(projectID, name, localFile string) {
	fr, _ := os.Open(localFile)
	defer fr.Close()

	r := bufio.NewReader(os.Stdin)
	ft := rex.FileTransformation{}
	ft.Position.Coordinates = make([]float64, 3)

	v := getInput("Pos X [m]: ", r)
	ft.Position.Coordinates[0] = interpretFloatInput(v, 0.0)
	v = getInput("Pos Y [m]: ", r)
	ft.Position.Coordinates[1] = interpretFloatInput(v, 0.0)
	v = getInput("Pos Z [m]: ", r)
	ft.Position.Coordinates[2] = interpretFloatInput(v, 0.0)
	v = getInput("Rot X [degrees]: ", r)
	ft.Rotation.X = interpretFloatInput(v, 0.0)
	v = getInput("Rot Y [degrees]: ", r)
	ft.Rotation.Y = interpretFloatInput(v, 0.0)
	v = getInput("Rot Z [degrees]: ", r)
	ft.Rotation.Z = interpretFloatInput(v, 0.0)
	v = getInput("Scale [default=1]: ", r)
	ft.Scale = interpretFloatInput(v, 1.0)
	err := rex.UploadProjectFile(RxConfig.Client, projectID, name, localFile, &ft, fr)
	console(err, "Success!")
}

func createProjectInteractive() {
	r := bufio.NewReader(os.Stdin)

	name := getInput("Name", r)

	var address rex.ProjectAddress
	var absoluteTransformation rex.ProjectTransformation

	address.AddressLine1 = getInput("AddressLine1", r)
	address.AddressLine2 = getInput("AddressLine2", r)
	address.AddressLine3 = getInput("AddressLine3", r)
	address.AddressLine4 = getInput("AddressLine4", r)
	address.PostCode = getInput("Postcode", r)
	address.City = getInput("City", r)
	address.Region = getInput("Region", r)
	address.Country = getInput("Country", r)

	lat, lon := GetGeoLocation(&address)
	absoluteTransformation.Position.Coordinates = make([]float64, 3)

	latInput := getInput(fmt.Sprintf("Lat [%f]", lat), r)
	absoluteTransformation.Position.Coordinates[1] = interpretFloatInput(latInput, lat)

	lonInput := getInput(fmt.Sprintf("Lon [%f]", lon), r)
	absoluteTransformation.Position.Coordinates[0] = interpretFloatInput(lonInput, lon)

	heightInput := getInput("Height [m]", r)
	absoluteTransformation.Position.Coordinates[2] = interpretFloatInput(heightInput, 0.0)

	absoluteTransformation.Rotation.X = 0.0
	absoluteTransformation.Rotation.Y = 0.0

	northingInput := getInput("Northing (0=north, 90=east)", r)
	absoluteTransformation.Rotation.Z = interpretFloatInput(northingInput, 0.0)

	fmt.Printf("\n\nProject %s ", name)
	err := rex.CreateProject(RxConfig.Client, RxConfig.Client.User.UserID, name, &address, &absoluteTransformation)
	console(err, "added successfully!")

	mapLink := openStreetMap
	mapLink += fmt.Sprintf("?mlat=%f", absoluteTransformation.Position.Coordinates[1])
	mapLink += fmt.Sprintf("&mlon=%f", absoluteTransformation.Position.Coordinates[0])
	mapLink += fmt.Sprintf("#map=17/%f", absoluteTransformation.Position.Coordinates[1])
	mapLink += fmt.Sprintf("/%f", absoluteTransformation.Position.Coordinates[0])
	mapLink += fmt.Sprintf("&layers=N")
	fmt.Printf("\n\nMap link: %s\n\n", mapLink)
}
