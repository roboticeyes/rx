/*
 * Author: Bernhard Reitinger
 * Date  : 2018
 */

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/breiting/rex"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RxConfig stores the values from the configuration file
var RxConfig = struct {

	// ClientID is the secret ID for authentication
	ClientID string

	// ClientSecret is the secret password for authentication
	ClientSecret string

	// This is the token which is temporarily stored in the config file
	// if it is expired then the client needs to re-authorize
	ClientToken string

	// The authenticated client for all REX operations
	AuthClient *rex.Client
}{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rx",
	Short: "REX command line tool for accessing the REX cloud API",
	Long: `
--------------------------------------------------------------
                          rx - (c) 2018
--------------------------------------------------------------

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rx.yml)")

	rootCmd.PersistentFlags().StringVar(&rex.RexBaseURL, "BaseURL", "https://rex.robotic-eyes.com", "REX cloud base URL")
	rootCmd.PersistentFlags().StringVar(&RxConfig.ClientID, "ClientID", "", "client id for the user")
	rootCmd.PersistentFlags().StringVar(&RxConfig.ClientSecret, "ClientSecret", "", "client secret for the user")

	viper.BindPFlag("ClientID", rootCmd.PersistentFlags().Lookup("ClientID"))
	viper.BindPFlag("ClientSecret", rootCmd.PersistentFlags().Lookup("ClientSecret"))
	viper.BindPFlag("BaseURL", rootCmd.PersistentFlags().Lookup("BaseURL"))
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

		viper.AddConfigPath(home)
		viper.SetConfigName(".rx")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
	} else {
		return
	}
	RxConfig.ClientID = viper.GetString("ClientID")
	RxConfig.ClientSecret = viper.GetString("ClientSecret")
	RxConfig.ClientToken = viper.GetString("AuthToken")
	rex.RexBaseURL = viper.GetString("BaseURL")

	// Authenticate user
	var err error
	RxConfig.AuthClient, err = rex.NewClient(RxConfig.ClientID, RxConfig.ClientSecret, nil)
	if err != nil {
		log.Fatal("Cannot get authentication token, please check your client credentials")
	}
}
