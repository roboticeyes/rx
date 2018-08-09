package cmd

import (
	"fmt"
	"os"

	"github.com/breiting/rex"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// ClientID is the secret ID for authentication
var ClientID string

// ClientSecret is the secret password for authentication
var ClientSecret string

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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rexcli.yml)")

	rootCmd.PersistentFlags().StringVar(&rex.RexBaseURL, "BaseURL", "https://rex.robotic-eyes.com", "REX cloud base URL")
	rootCmd.PersistentFlags().StringVar(&ClientID, "ClientID", "", "client id for the user")
	rootCmd.PersistentFlags().StringVar(&ClientSecret, "ClientSecret", "", "client secret for the user")

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
		viper.SetConfigName(".rexcli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
	} else {
		panic(err)
	}
	ClientID = viper.GetString("ClientID")
	ClientSecret = viper.GetString("ClientSecret")
	rex.RexBaseURL = viper.GetString("BaseURL")
}
