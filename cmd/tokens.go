package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:   "tokens",
	Short: "use for generating tokens",
	Long: `This command allows for generation for various tokens for authentication.
Tokens can be generated for user activation, autnetication and password resets`,
}

var tokenActivationCmd = &cobra.Command{
	Use:   "activation [email]",
	Short: "generate a activation token for user",
	Long: `Generate a new activation for a user who is not activated.
A user may be able to request a new token in case the current token in exipired
or they not recevie the intial activation token for some reason`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestToken("/v1/tokens/activation", cmd, args)
	},
}

var tokenPasswordResetCmd = &cobra.Command{
	Use:   "password-reset [email]",
	Short: "generate a password reset token",
	Long:  `Generate a new password reset token for user in case they need to reset their password`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return requestToken("/v1/tokens/password-reset", cmd, args)
	},
}

func requestToken(url string, cmd *cobra.Command, args []string) error {
	var email string
	flag := cmd.Flag("email")
	if len(args) < 1 && !flag.Changed {
		return fmt.Errorf("email is required in order to request a activation token")
	}

	if len(args) > 0 {
		email = args[0]
	} else {
		email = userEmail
	}

	jsonMap := newJSONMap()
	jsonMap.add("email", email)

	jsReader, err := jsonMap.createJSONReader()
	if err != nil {
		return err
	}

	resp := apiClient.NewRequest(http.MethodPost, url, jsReader, nil)
	if resp.Err != nil {
		return resp.Err
	}

	if resp.Code != http.StatusAccepted {
		return customError(cmd, resp.Body)
	}

	fmt.Fprintln(os.Stdout, resp.Body)
	return nil
}

func init() {
	tokenActivationCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email of registered user")

	tokenPasswordResetCmd.Flags().StringVarP(&userEmail, "email", "e", "", "email for password reset")

	tokenCmd.AddCommand(tokenActivationCmd)
	tokenCmd.AddCommand(tokenPasswordResetCmd)
}
