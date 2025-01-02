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

var userRegisterCmd = &cobra.Command{
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

		resp := apiClient.NewRequest(http.MethodPost, "/v1/users", jsReader, nil)
		if resp.Err != nil {
			return resp.Err
		}

		if resp.Code != http.StatusAccepted {
			return customError(cmd, resp.Body)
		}

		fmt.Fprintf(os.Stdout, "User Created!\nUser Details:\n%s\n", resp.Body)
		return nil
	},
}

var userActivateCmd = &cobra.Command{
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

		resp := apiClient.NewRequest(http.MethodPut, "/v1/users/activated", js, nil)
		if resp.Err != nil {
			return resp.Err
		}

		if resp.Code != http.StatusOK {
			return customError(cmd, resp.Body)
		}

		fmt.Fprintf(os.Stdout, "User Activated!\n%s\n", resp.Body)

		return nil
	},
}

var userPasswordResetCmd = &cobra.Command{
	Use:   "password-reset",
	Short: "reset password with a token",
	Long:  "reset a users password with a token as sent to their email",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonMap := newJSONMap()
		jsonMap.add("token", passwordResetToken)
		jsonMap.add("password", userPassword)

		js, err := jsonMap.createJSONReader()
		if err != nil {
			return err
		}

		resp := apiClient.NewRequest(http.MethodPut, "/v1/users/password", js, nil)
		if resp.Err != nil {
			return resp.Err
		}

		if resp.Code != http.StatusOK {
			return customError(cmd, resp.Body)
		}

		fmt.Fprint(os.Stdout, resp.Body)
		return nil
	},
}

func init() {
	userRegisterCmd.Flags().StringVarP(&userName, "name", "n", "", "name of user to register")
	userRegisterCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of user to register")
	userRegisterCmd.Flags().StringVarP(&userPassword, "password", "p", "", "password of user to register")
	userRegisterCmd.MarkFlagRequired("name")
	userRegisterCmd.MarkFlagRequired("email")
	userRegisterCmd.MarkFlagRequired("password")

	userActivateCmd.Flags().StringVarP(&activationToken, "token", "t", "", "activation token for the user")

	userPasswordResetCmd.Flags().StringVarP(&userPassword, "password", "p", "", "new password for the user")
	userPasswordResetCmd.Flags().StringVarP(&passwordResetToken, "token", "t", "", "token for password reset")
	userPasswordResetCmd.MarkFlagRequired("password")
	userPasswordResetCmd.MarkFlagRequired("token")

	userCmd.AddCommand(userRegisterCmd)
	userCmd.AddCommand(userActivateCmd)
	userCmd.AddCommand(userPasswordResetCmd)
}
