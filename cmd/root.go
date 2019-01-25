package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/breiting/rex"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var version = "v0.1.0"

var cfgFile string

// RxConfig stores the values from the configuration file
var RxConfig = struct {
	ClientID     string      // ClientID is the secret ID for authentication
	ClientSecret string      // ClientSecret is the secret password for authentication
	Client       *rex.Client // The authenticated client for all REX operations
}{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rx",
	Short: "REX command line tool for accessing the REX cloud API",
	Long: getCmdLineHeader() + `
rx is a command line tool for accessing the REX cloud API.
For further information please see our support page:

                https://support.robotic-eyes.com

`}

// Execute the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/rx/config.json)")
	rootCmd.PersistentFlags().StringVar(&rex.RexBaseURL, "BaseURL", "https://rex.robotic-eyes.com", "REX cloud base URL")

	rootCmd.PersistentFlags().StringVar(&RxConfig.ClientID, "ClientID", "", "client id for the user (required for login)")
	rootCmd.PersistentFlags().StringVar(&RxConfig.ClientSecret, "ClientSecret", "", "client secret for the user (required for login)")

}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		rxConfigPath := path.Join(home, ".config/rx")
		viper.AddConfigPath(rxConfigPath)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Printf("Found configuration file: %s\n", viper.ConfigFileUsed())

		if RxConfig.ClientID == "" {
			RxConfig.ClientID = viper.GetString("ClientID")
		}
		if RxConfig.ClientSecret == "" {
			RxConfig.ClientSecret = viper.GetString("ClientSecret")
		}
		if viper.IsSet("BaseURL") {
			rex.RexBaseURL = viper.GetString("BaseURL")
		}
	}

	// Authenticate user (either use token file or if not working ClientId and ClientSecret)
	const tokenFile = "token"
	buf, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		fmt.Println("No token found, try to use ClientID and ClientSecret for authentication")

		fmt.Println("ClientID:     ", RxConfig.ClientID)
		fmt.Println("ClientSecret: ", RxConfig.ClientSecret)

		RxConfig.Client = rex.NewClient(nil)
		err := RxConfig.Client.Login(RxConfig.ClientID, RxConfig.ClientSecret)
		if err != nil {
			fmt.Println("Cannot login, please check your client credentials")
			os.Exit(1)
		}

		buf, err := json.Marshal(&RxConfig.Client.Token)
		err = ioutil.WriteFile(tokenFile, buf, 0600)
		if err != nil {
			fmt.Println("Cannot write token file")
			os.Exit(1)
		}
		fmt.Println("Successfully got token and stored it in file: ", tokenFile)
		return
	}
	var token oauth2.Token
	err = json.Unmarshal(buf, &token)
	RxConfig.Client, err = rex.NewClientWithToken(token, nil)
	if err != nil {
		fmt.Println("Cannot login with token, got error: ", err, ". Please use login first.")
		os.Exit(1)
	}
}
