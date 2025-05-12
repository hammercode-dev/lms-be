package migrate

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use: "create [name]",
	Short: "create migration file",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dirName := "database/migrations"
		name := strings.ToLower(args[0])
		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s/%s_%s.sql", dirName, timestamp, name)
		content := "-- +migrate Up\n\n-- +migrate Down\n"

		err := os.MkdirAll(dirName, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating migrations folder:", err)
			return
		}

		err = os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			fmt.Println("Error creating migration file:", err)
			return
		}

		fmt.Println("Migration created:", filename)
	},
}