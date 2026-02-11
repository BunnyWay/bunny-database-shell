package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/libsql/libsql-shell-go/pkg/shell"
	"github.com/libsql/libsql-shell-go/pkg/shell/enums"
	"github.com/spf13/cobra"
)

var (
	flagURL       string
	flagAuthToken string
)

func main() {
	godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "bunny-database-shell",
		Short: "Connect to a Bunny Database shell",
		RunE: func(c *cobra.Command, args []string) error {
			url := resolve(flagURL, "BUNNY_DB_URL", "Database URL: ")
			authToken := resolve(flagAuthToken, "BUNNY_DB_TOKEN", "Auth Token: ")

			if url == "" {
				return fmt.Errorf("database URL is required")
			}

			welcomeMessage := "Bunny Database Shell\n\n"

			config := shell.ShellConfig{
				DbUri:          url,
				AuthToken:      authToken,
				InF:            os.Stdin,
				OutF:           os.Stdout,
				ErrF:           os.Stderr,
				HistoryMode:    enums.PerDatabaseHistory,
				HistoryName:    "bunny-database-shell",
				WelcomeMessage: &welcomeMessage,
			}

			return shell.RunShell(config)
		},
	}

	rootCmd.Flags().StringVar(&flagURL, "url", "", "Bunny Database URL (libsql:// or wss://)")
	rootCmd.Flags().StringVar(&flagAuthToken, "auth-token", "", "Authentication token")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func resolve(flag, envKey, promptLabel string) string {
	if flag != "" {
		return flag
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	return prompt(promptLabel)
}

func prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
