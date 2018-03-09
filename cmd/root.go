package cmd

import (
	"fmt"
	"os"

	"github.com/pivotal-cloudops/omen/internal/opsman"
	"github.com/pivotal-cloudops/omen/internal/userio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ENV_OPSMAN_HOST          = "OPSMAN_HOSTNAME"
	ENV_OPSMAN_USERNAME      = "OPSMAN_USER"
	ENV_OPSMAN_PASSWORD      = "OPSMAN_PASSWORD"
	ENV_OPSMAN_CLIENT_ID     = "OPSMAN_CLIENT_ID"
	ENV_OPSMAN_CLIENT_SECRET = "OPSMAN_CLIENT_SECRET"
)

var rp = userio.ReportPrinter{}

var rootCmd = &cobra.Command{
	Use:   "omen",
	Short: "omen is a phenomenal supplemental too to the Pivotal OM CLI",
	Long:  "omen adds functionality helpful to PCF operators",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No command to run, use --help for a list of available commands")
		os.Exit(1)
	},
}

func init() {
	var omHost, omUser, omPassword, omClientID, omClientSecret string

	rootCmd.PersistentFlags().StringVarP(&omHost, "om_host", "H", "",
		fmt.Sprintf("URL to Opsmanager (Defaults to Env Var $%s)", ENV_OPSMAN_HOST))

	rootCmd.PersistentFlags().StringVarP(&omUser, "om_user", "u", "",
		fmt.Sprintf("Opsmanager User Name (Defaults to Env Var $%s)", ENV_OPSMAN_USERNAME))

	rootCmd.PersistentFlags().StringVarP(&omPassword, "om_password", "p", "",
		fmt.Sprintf("Opsmanager Password (Defaults to Env Var $%s)", ENV_OPSMAN_PASSWORD))

	rootCmd.PersistentFlags().StringVarP(&omClientID, "om_client_id", "c", "",
		fmt.Sprintf("Opsmanager Client ID (Defaults to Env Var $%s)", ENV_OPSMAN_CLIENT_ID))

	rootCmd.PersistentFlags().StringVarP(&omClientSecret, "om_client_secret", "s", "",
		fmt.Sprintf("Opsmanager Client Secret (Defaults to Env Var $%s)", ENV_OPSMAN_CLIENT_SECRET))

	viper.BindPFlag("omHost", rootCmd.PersistentFlags().Lookup("om_host"))
	viper.BindEnv("omHost", ENV_OPSMAN_HOST)

	viper.BindPFlag("omUser", rootCmd.PersistentFlags().Lookup("om_user"))
	viper.BindEnv("omUser", ENV_OPSMAN_USERNAME)

	viper.BindPFlag("omPassword", rootCmd.PersistentFlags().Lookup("om_password"))
	viper.BindEnv("omPassword", ENV_OPSMAN_PASSWORD)

	viper.BindPFlag("omClientID", rootCmd.PersistentFlags().Lookup("om_client_id"))
	viper.BindEnv("omClientID", ENV_OPSMAN_CLIENT_ID)

	viper.BindPFlag("omClientSecret", rootCmd.PersistentFlags().Lookup("om_client_secret"))
	viper.BindEnv("omClientSecret", ENV_OPSMAN_CLIENT_SECRET)

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(diagnosticsCmd)
	rootCmd.AddCommand(manifestsCmd)
	rootCmd.AddCommand(applyChangesCmd)
	rootCmd.AddCommand(stagedTilesCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getOpsmanClient() opsman.Client {
	url := viper.GetString("omHost")
	user := ""
	secret := ""
	clientID := viper.GetString("omClientID")
	clientSecret := viper.GetString("omClientSecret")

	if url == "" {
		fmt.Println("Opsman host is required. Please specify by flag or environment variable")
		os.Exit(1)
	}

	if clientID == "" && clientSecret == "" {
		user = viper.GetString("omUser")
		secret = viper.GetString("omPassword")

		if user == "" {
			fmt.Println("Opsman user is required. Please specify by flag or environment variable")
			os.Exit(1)
		}

		if secret == "" {
			fmt.Println("Opsman user secret is required. Please specify by flag or environment variable")
			os.Exit(1)
		}
	} else {
		if clientID == "" {
			fmt.Println("Opsman client ID is required. Please specify by flag or environment variable")
			os.Exit(1)
		}

		if clientSecret == "" {
			fmt.Println("Opsman client secret is required. Please specify by flag or environment variable")
			os.Exit(1)
		}
	}

	return opsman.NewClient(url, user, secret, clientID, clientSecret)
}
