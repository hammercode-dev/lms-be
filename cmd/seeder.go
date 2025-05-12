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

		gooseCmd := exec.Command("goose", "create", "seed_"+args[0], "-dir", seederDir, "sql")
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose create migration error: ", err)
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
		gooseCmd := exec.Command("goose", "up", "-dir", seederDir, cfg.DB_POSTGRES_DSN)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose up migration error: ", err)
			return
		}

	},
}

var seedDown = &cobra.Command{
	Use:   "seed:down",
	Short: "seed down",
	Long:  "seed down",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()

		gooseCmd := exec.Command("goose", "down", "-dir", seederDir, cfg.DB_POSTGRES_DSN)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose down migration error: ", err)
			return
		}

	},
}
