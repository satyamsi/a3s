package authcmd

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/manipulate/manipcli"
)

func makeCheckCmd(mmaker manipcli.ManipulatorMaker) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "check",
		Aliases:          []string{"verify"},
		Short:            "Check the token",
		TraverseChildren: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if err := HandleAutoAuth(
				mmaker,
				"",
				nil,
				nil,
				false,
			); err != nil {
				return fmt.Errorf("auto auth error: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			fToken := viper.GetString("token")
			fPrint := viper.GetBool("print")
			fQRCode := viper.GetBool("qrcode")

			return CheckToken(fToken, fPrint, fQRCode)
		},
	}

	cmd.Flags().Bool("print", false, "Print the token string.")

	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cmd.Flags().MarkHidden("namespace")
		cmd.Flags().MarkHidden("audience")
		cmd.Flags().MarkHidden("cloak")
		cmd.Flags().MarkHidden("validity")
		cmd.Flags().MarkHidden("encoding")
		cmd.Flags().MarkHidden("restrict-namespace")
		cmd.Flags().MarkHidden("restrict-permissions")
		cmd.Flags().MarkHidden("restrict-network")
		cmd.Parent().HelpFunc()(cmd, args)
	})

	return cmd
}

func CheckToken(token string, printRaw bool, qrcode bool) error {
	claims := jwt.MapClaims{}
	p := jwt.Parser{}

	t, _, err := p.ParseUnverified(token, &claims)
	if err != nil {
		return err
	}

	data, err := prettyjson.Marshal(claims)
	if err != nil {
		return err
	}

	fmt.Println("alg:", t.Method.Alg())
	fmt.Println("kid:", t.Header["kid"])
	fmt.Println()

	fmt.Println(string(data))

	if printRaw {
		fmt.Println()
		printToken(token, qrcode)
	}

	return nil
}
