package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/libsql/libsql-shell-go/pkg/shell"
	"github.com/libsql/libsql-shell-go/pkg/shell/enums"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	flagURL       string
	flagAuthToken string
	flagExec      string
)

func main() {
	godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "bunny-database-shell [sql]",
		Short: "Connect to a Bunny Database shell",
		Args:  cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			url := resolve(flagURL, "BUNNY_DB_URL", "Database URL: ", false)
			authToken := resolve(flagAuthToken, "BUNNY_DB_TOKEN", "Auth Token: ", true)

			if url == "" {
				return fmt.Errorf("database URL is required")
			}

			config := shell.ShellConfig{
				DbUri:       url,
				AuthToken:   authToken,
				InF:         os.Stdin,
				OutF:        os.Stdout,
				ErrF:        os.Stderr,
				HistoryMode: enums.PerDatabaseHistory,
				HistoryName: "bunny-database-shell",
			}

			stmt := flagExec
			if stmt == "" && len(args) > 0 {
				stmt = strings.Join(args, " ")
			}
			if stmt == ".dump" {
				return dump(url, authToken)
			}
			if stmt != "" {
				return shell.RunShellLine(config, stmt)
			}

			welcomeMessage := "\nYou are connected to your Bunny Database shell.\n\nType \".quit\" to exit, \".help\" for commands.\n\n"
			config.WelcomeMessage = &welcomeMessage
			return shell.RunShell(config)
		},
	}

	rootCmd.Flags().StringVar(&flagURL, "url", "", "Bunny Database URL (libsql:// or wss://)")
	rootCmd.Flags().StringVar(&flagAuthToken, "auth-token", "", "Authentication token")
	rootCmd.Flags().StringVarP(&flagExec, "exec", "e", "", "Execute a SQL statement and exit")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func resolve(flag, envKey, promptLabel string, secret bool) string {
	if flag != "" {
		return flag
	}
	if v := os.Getenv(envKey); v != "" {
		return v
	}
	if secret {
		return promptSecret(promptLabel)
	}
	return prompt(promptLabel)
}

func prompt(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptSecret(label string) string {
	fmt.Print(label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func dumpURL(dbURL string) string {
	switch {
	case strings.HasPrefix(dbURL, "libsql://"):
		return strings.Replace(dbURL, "libsql://", "https://", 1)
	case strings.HasPrefix(dbURL, "wss://"):
		return strings.Replace(dbURL, "wss://", "https://", 1)
	case strings.HasPrefix(dbURL, "ws://"):
		return strings.Replace(dbURL, "ws://", "http://", 1)
	default:
		return dbURL
	}
}

func dump(dbURL, authToken string) error {
	req, err := http.NewRequest("GET", dumpURL(dbURL)+"/dump", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+authToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("dump failed (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Print(line)
		if err == io.EOF {
			return nil
		}
	}
}
