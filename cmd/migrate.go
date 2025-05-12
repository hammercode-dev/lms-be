package cmd

import (
	"os"
	"os/exec"

	"github.com/hammer-code/lms-be/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//	var dbMigrate = &cobra.Command{
//		Use:   "migrate",
//		Short: "migrate database",
//		Long:  "migrate database",
//		Run: func(cmd *cobra.Command, args []string) {
//			// fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
//			// ctx := context.Background()
//
//			cfg := config.GetConfig()
//
//			db := config.GetDatabase(postgres.Dialector{
//				Config: &postgres.Config{
//					DSN: cfg.DB_POSTGRES_DSN,
//				}})
//
//			err := db.AutoMigrate(&domain.User{}, &domain.LogoutToken{})
//			if err != nil {
//				logrus.Error(err)
//				return
//			}
//
//			logrus.Info("Migrated")
//		},
//	}
var (
	migrationDir = "database/migration"
)

var createMigration = &cobra.Command{
	Use:   "migrate:create",
	Short: "migrate create",
	Long:  "migrate create",
	Run: func(cmd *cobra.Command, args []string) {

		gooseCmd := exec.Command("goose", "create", "table"+args[0], "-dir", migrationDir, "sql")
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose create migration error: ", err)
			return
		}

	},
}

var migrateUp = &cobra.Command{
	Use:   "migrate:up",
	Short: "migrate up",
	Long:  "migrate up",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()
		gooseCmd := exec.Command("goose", "up", "-dir", migrationDir, cfg.DB_POSTGRES_DSN)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose up migration error: ", err)
			return
		}

	},
}

var migrateDown = &cobra.Command{
	Use:   "migrate:down",
	Short: "migrate down",
	Long:  "migrate down",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.GetConfig()

		gooseCmd := exec.Command("goose", "down", "-dir", migrationDir, cfg.DB_POSTGRES_DSN)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose down migration error: ", err)
			return
		}

	},
}
