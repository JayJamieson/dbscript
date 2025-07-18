package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/JayJamieson/dbscript/pkg/mysql"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	user     string
	host     string
	port     int
	password string
	schema   string
	tables   []string
	handler  string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start CDC event processing",
	Long:  `Start processing CDC events from MySQL with JavaScript handlers.`,
	Run: func(cmd *cobra.Command, args []string) {

		if password == "" {
			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
				os.Exit(1)
			}
			password = string(bytePassword)
			fmt.Println() // Add newline after password input
		} else {
			fmt.Fprintf(os.Stderr, "Warning: Using plain text password from command line is not secure\n")
		}

		listener, err := mysql.NewBinlogListener(&mysql.BinlogListenerOptions{
			Host:     host,
			Port:     port,
			User:     user,
			Schema:   schema,
			Tables:   tables,
			Password: password,
		})

		if err != nil {
			listener.Logger.Error("Error creating BinlogListener", "error", err)
			os.Exit(1)
		}

		listener.Logger.Info("Starting CDC processing with:",
			slog.Group("config", slog.String("schema", schema),
				slog.String("host", host),
				slog.Int("port", port),
				slog.String("user", user),
				slog.String("tables", strings.Join(tables, ",")),
				slog.String("handler", handler)),
		)

		sig := make(chan os.Signal, 1)

		signal.Notify(sig, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		go func() {
			if err := listener.Listen(); err != nil {
				listener.Logger.Error("Error starting dbscript", "error", err)
			}
		}()

		<-sig
		listener.Close()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&user, "user", "u", "", "Database user")
	startCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Database host")
	startCmd.Flags().IntVarP(&port, "port", "p", 3306, "Database port")
	startCmd.Flags().StringVar(&password, "password", "", "Database password (leave empty to prompt)")
	startCmd.Flags().StringVar(&schema, "schema", "", "Database schema name")
	startCmd.Flags().StringSliceVar(&tables, "tables", []string{}, "Tables to monitor for changes")
	startCmd.Flags().StringVar(&handler, "handler", "", "JavaScript handler file")

	startCmd.MarkFlagRequired("user")
	startCmd.MarkFlagRequired("schema")
	startCmd.MarkFlagRequired("tables")
	startCmd.MarkFlagRequired("handler")
}
