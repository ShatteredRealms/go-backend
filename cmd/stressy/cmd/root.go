package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/ShatteredRealms/go-backend/cmd/stressy/internal"
	"github.com/ShatteredRealms/go-backend/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	conf    *config.GlobalConfig
	showCfg bool
	stop    context.CancelFunc
	k       *internal.KeycloakManager

	numAdmins int
	numUsers  int

	rootCmd = &cobra.Command{
		Use:   "stressy",
		Short: "Stress tests SRO applications",
		Long: `A stress test cli tool for stress testing Shattered Realms Online.

Tasks:
* Create temporary users of all authroziation levels
* Make valid, and invalid, requests to specified SRO services`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			ctx, stop = signal.NotifyContext(ctx, os.Interrupt)
			defer stop()

			conf = config.NewGlobalConfig(context.Background())
			b, err := json.MarshalIndent(conf, "", " ")
			if err != nil {
				fmt.Print("Error: unable to decode config")
				os.Exit(1)
			}

			if showCfg {
				fmt.Printf("Using config:\n%s", string(b))
			}

			go func() {
				k = internal.NewKeycloakManager(ctx, conf, numAdmins, numUsers, 1)
				fmt.Printf("Setting up test accounts: %d admins and %d users\n", numAdmins, numUsers)
				err = k.Setup(ctx)
				if err != nil {
					fmt.Printf("error setting up keycloak: %v\n", err)
					return
				}
				fmt.Println("Finished up test accounts")
			}()

			<-ctx.Done()
			shutdown()
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/stresy/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&showCfg, "showConfig", "c", false, "show the SRO config being used")
	rootCmd.PersistentFlags().IntVarP(&numAdmins, "admins", "a", 2, "number of admins to simulate")
	rootCmd.PersistentFlags().IntVarP(&numUsers, "users", "u", 50, "number of users to simulate")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home + "/.config/stressy")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func shutdown() {
	fmt.Println("\nExiting...")

	if k != nil {
		fmt.Println("Removing temp accounts")
		err := k.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("Errors encountered shutting down: %v", err)
		} else {
			fmt.Println("Temp accounts removed")
		}
	}

	fmt.Println("Finish exiting")
	stop()
}
