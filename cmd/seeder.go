package cmd

import (
	"os"
	"os/exec"

	"github.com/hammer-code/lms-be/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	seederDir = "database/seeder"
)

var createSeeder = &cobra.Command{
	Use:   "seed:create",
	Short: "seed create",
	Long:  "seed create",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			logrus.Error("Seeder name is required, e.g. `go run main.go seed:create users`")
			return
		}
		gooseCmd := exec.Command("goose", "create", "seed_"+args[0], "-dir", seederDir, "sql")
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose create seed error: ", err)
			return
		}

	},
}

var seedUp = &cobra.Command{
	Use:   "seed:up",
	Short: "seed up",
	Long:  "seed up",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()
		gooseCmd := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "up", "-dir", seederDir)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose up seed error: ", err)
			return
		}

	},
}

var seedStatus = &cobra.Command{
	Use:   "seed:status",
	Short: "show seed status",
	Long:  "show status of goose seed migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		gooseCmd := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "status", "-dir", seederDir)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose seed status error: ", err)
			return
		}
	},
}
var seedReset = &cobra.Command{
	Use:   "seed:reset",
	Short: "seed reset (down + up)",
	Long:  "reset all seeders",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		// Down
		down := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "down", "-dir", seederDir)
		down.Stdout = os.Stdout
		down.Stderr = os.Stderr
		if err := down.Run(); err != nil {
			logrus.Error("goose seed down error: ", err)
			return
		}

		// Up
		up := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "up", "-dir", seederDir)
		up.Stdout = os.Stdout
		up.Stderr = os.Stderr
		if err := up.Run(); err != nil {
			logrus.Error("goose seed up error: ", err)
			return
		}
	},
}
