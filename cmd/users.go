package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tiwanakd/greenlight-cli/client"
)

var userCmd = &cobra.Command{
	Use:   "users",
	Short: "make requests related to a user",
	Long: `This command provies to make user related requests
a new user can register, activate and update their password`,
}

var registerUserCmd = &cobra.Command{
	Use:   "register",
	Short: "register a new user",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMap := client.NewJSONMap()
		jsonMap.Add("name", userName)
		jsonMap.Add("email", userEmail)
		jsonMap.Add("password", userPassword)

		jsReader, err := jsonMap.CreateJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.PostRequest("/v1/users", jsReader)
		if err != nil {
			return err
		}

		if code != http.StatusAccepted {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return errors.New(body)
		}

		fmt.Fprintf(os.Stdout, "User Created!\nUser Details:\n%s\n", body)
		return nil
	},
}

func init() {
	registerUserCmd.Flags().StringVarP(&userName, "name", "n", "", "name of user to register")
	registerUserCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of user to register")
	registerUserCmd.Flags().StringVarP(&userPassword, "password", "p", "", "password of user to register")
	registerUserCmd.MarkFlagRequired("name")
	registerUserCmd.MarkFlagRequired("email")
	registerUserCmd.MarkFlagRequired("password")

	userCmd.AddCommand(registerUserCmd)
}
