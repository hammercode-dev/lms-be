package migrate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hammer-code/lms-be/config"
	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		conn, err := pgx.Connect(context.Background(), cfg.DB_POSTGRES_DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		files, _ := filepath.Glob("database/migrations/*.sql")
		for _, file := range files {
			content, _ := os.ReadFile(file)
			parts := strings.Split(string(content), "-- +migrate Down")
			upSQL := strings.Replace(parts[0], "-- +migrate Up", "", 1)

			if _, err := conn.Exec(context.Background(), upSQL); err != nil {
				fmt.Println("Failed to run migration:", file, err)
				return
			}

			fmt.Println("Migrated:", file)
		}
	},
}