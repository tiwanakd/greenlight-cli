package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiwanakd/greenlight-cli/client"
)

var (
	userName     string
	userEmail    string
	userPassword string

	activationToken    string
	authorizationToken string
	passwordResetToken string
)

var apiClient = client.New()

var rootCmd = &cobra.Command{
	Use:   "greenlight",
	Short: "greenlight is client to interact with greenlight-api",
	Long: `This is client that is build to interact with greenlight-api project
the api is deployed at greenlight-api.com this cli-clent provides various 
commands that allows to interract with this api`,
}

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "check the health of the greenlight api",
	RunE: func(cmd *cobra.Command, args []string) error {
		err, code, body := apiClient.NewRequest(http.MethodGet, "/v1/healthcheck", http.NoBody, nil)
		if err != nil {
			return err
		}

		if code != http.StatusOK {
			return customError(cmd, body)
		}

		fmt.Printf("Body: %s\n", body)
		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login to authenticate",
	Long:  "login for the current terminal session (email and password required)",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMap := newJSONMap()
		jsonMap.add("email", userEmail)
		jsonMap.add("password", userPassword)

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		err, code, body := apiClient.NewRequest(http.MethodPost, "/v1/tokens/authentication", js, nil)
		if err != nil {
			return err
		}

		if code != http.StatusCreated {
			return customError(cmd, body)
		}

		err = addAuthTokenToFile([]byte(body))
		if err != nil {
			return err
		}

		fmt.Println("Logged in Successfully!")
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "logout the current user",
	Long:  "removes the token file and which works as logging off the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := removeAuth()
		if err != nil {
			return err
		}

		fmt.Println("You have been logged out successfully!")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of user to register")
	loginCmd.Flags().StringVarP(&userPassword, "password", "p", "", "password of user to register")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(movieCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(healthCheckCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
