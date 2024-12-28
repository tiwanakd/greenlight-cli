package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
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
		jsonMap := newJSONMap()
		jsonMap.add("name", userName)
		jsonMap.add("email", userEmail)
		jsonMap.add("password", userPassword)

		jsReader, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.NewRequest(http.MethodPost, "/v1/users", jsReader, nil)
		if err != nil {
			return err
		}

		if code != http.StatusAccepted {
			return customError(cmd, body)
		}

		fmt.Fprintf(os.Stdout, "User Created!\nUser Details:\n%s\n", body)
		return nil
	},
}

var activateUserCmd = &cobra.Command{
	Use:   "activate",
	Short: "active a user with provided token",
	Long: `this command will active the user with the provided token; 
given the token is valid anod not expired`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string
		tokenFlag := cmd.Flag("token")
		if len(args) < 1 && !tokenFlag.Changed {
			return fmt.Errorf("token is required for user activation")
		}

		if len(args) > 0 {
			token = args[0]
		} else {
			token = activationToken
		}

		jsonMap := newJSONMap()
		jsonMap.add("token", token)

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		err, code, _, body := apiClient.NewRequest(http.MethodPut, "/v1/users/activated", js, nil)
		if err != nil {
			return err
		}

		if code != http.StatusOK {
			return customError(cmd, body)
		}

		fmt.Printf("User Activated!\n%s\n", body)
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

	activateUserCmd.Flags().StringVarP(&activationToken, "token", "t", "", "activation token for the user")

	userCmd.AddCommand(registerUserCmd)
	userCmd.AddCommand(activateUserCmd)
}
