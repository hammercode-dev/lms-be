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

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Revert all down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		conn, err := pgx.Connect(context.Background(), cfg.DB_POSTGRES_DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		files, _ := filepath.Glob("database/migrations/*.sql")
		// reverse order
		for i := len(files) - 1; i >= 0; i-- {
			file := files[i]
			content, _ := os.ReadFile(file)
			parts := strings.Split(string(content), "-- +migrate Down")
			if len(parts) < 2 {
				continue
			}
			downSQL := parts[1]

			if _, err := conn.Exec(context.Background(), downSQL); err != nil {
				fmt.Println("Failed to rollback migration:", file, err)
				return
			}

			fmt.Println("Rolled back:", file)
		}
	},
}