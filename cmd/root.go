package cmd

import (
	"fmt"
	"os"

	"github.com/pivotal-cloudops/omen/internal/opsman"
	"github.com/pivotal-cloudops/omen/internal/sessions"
	"github.com/pivotal-cloudops/omen/internal/userio"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envOpsmanHost         = "OPSMAN_HOSTNAME"
	envOpsmanUsername     = "OPSMAN_USER"
	envOpsmanPassword     = "OPSMAN_PASSWORD"
	envOpsmanClientId     = "OPSMAN_CLIENT_ID"
	envOpsmanClientSecret = "OPSMAN_CLIENT_SECRET"

	keyTarget       = "omTarget"
	keyUser         = "omUser"
	keyPassword     = "omPassword"
	keyClientId     = "omClientID"
	keyClientSecret = "omClientSecret"
	keyForceLogout  = "forceLogout"
)

var rp = userio.ReportPrinter{}
var tr = userio.NewTableReporter()

var rootCmd = &cobra.Command{
	Use:   "omen",
	Short: "omen is a phenomenal supplemental tool to the Pivotal OM CLI",
	Long:  "omen adds functionality helpful to PCF operators",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("No command to run, use --help for a list of available commands")
		os.Exit(1)
	},
}

func init() {
	var omHost, omUser, omPassword, omClientID, omClientSecret string
	var forceLogout bool

	rootCmd.PersistentFlags().StringVarP(&omHost, "target", "t", "",
		fmt.Sprintf("URL to Opsmanager (Defaults to Env Var $%s)", envOpsmanHost))

	rootCmd.PersistentFlags().StringVarP(&omUser, "username", "u", "",
		fmt.Sprintf("Opsmanager User Name (Defaults to Env Var $%s)", envOpsmanUsername))

	rootCmd.PersistentFlags().StringVarP(&omPassword, "password", "p", "",
		fmt.Sprintf("Opsmanager Password (Defaults to Env Var $%s)", envOpsmanPassword))

	rootCmd.PersistentFlags().StringVarP(&omClientID, "client-id", "c", "",
		fmt.Sprintf("Opsmanager Client ID (Defaults to Env Var $%s)", envOpsmanClientId))

	rootCmd.PersistentFlags().StringVarP(&omClientSecret, "client-secret", "s", "",
		fmt.Sprintf("Opsmanager Client Secret (Defaults to Env Var $%s)", envOpsmanClientSecret))

	rootCmd.PersistentFlags().BoolVarP(&forceLogout, "force-logout", "f", false,
		"(optional) Log all other users out of opsman before attempting action")

	viper.BindPFlag(keyTarget, rootCmd.PersistentFlags().Lookup("target"))
	viper.BindEnv(keyTarget, envOpsmanHost)

	viper.BindPFlag(keyUser, rootCmd.PersistentFlags().Lookup("username"))
	viper.BindEnv(keyUser, envOpsmanUsername)

	viper.BindPFlag(keyPassword, rootCmd.PersistentFlags().Lookup("password"))
	viper.BindEnv(keyPassword, envOpsmanPassword)

	viper.BindPFlag(keyClientId, rootCmd.PersistentFlags().Lookup("client-id"))
	viper.BindEnv(keyClientId, envOpsmanClientId)

	viper.BindPFlag(keyClientSecret, rootCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindEnv(keyClientSecret, envOpsmanClientSecret)

	viper.BindPFlag(keyForceLogout, rootCmd.PersistentFlags().Lookup("force-logout"))

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(diagnosticsCmd)
	rootCmd.AddCommand(manifestsCmd)
	rootCmd.AddCommand(applyChangesCmd)
	rootCmd.AddCommand(stagedTilesCmd)
	rootCmd.AddCommand(toggleErrandsCmd)
	rootCmd.AddCommand(stemcellUpdatesCmd)
	rootCmd.AddCommand(errandsCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupOpsmanClient() opsman.Client {
	url := viper.GetString(keyTarget)
	user := ""
	secret := ""
	clientID := viper.GetString(keyClientId)
	clientSecret := viper.GetString(keyClientSecret)
	forceLogout := viper.GetBool(keyForceLogout)

	if url == "" {
		fmt.Println("Opsman host is required. Please specify by flag or environment variable")
		os.Exit(1)
	}

	if clientID == "" && clientSecret == "" {
		user = viper.GetString(keyUser)
		secret = viper.GetString(keyPassword)

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

	client := opsman.NewClient(url, user, secret, clientID, clientSecret)

	if forceLogout == true {
		fmt.Println("Logging out all active opsman sessions.")
		err := sessions.NewSessionManager(client).ClearAll()
		if err != nil {
			fmt.Println("Failed to clear sessions with error: ", err.Error())
			os.Exit(1)
		}
	}

	return client
}
