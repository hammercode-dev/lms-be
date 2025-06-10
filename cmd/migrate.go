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

		gooseCmd := exec.Command("goose", "create", "table_"+args[0], "-dir", migrationDir, "sql")
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
		gooseCmd := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "up", "-dir", migrationDir)
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

		gooseCmd := exec.Command("goose", "postgres", cfg.DB_POSTGRES_DSN, "down", "-dir", migrationDir)
		gooseCmd.Stdout = os.Stdout
		gooseCmd.Stderr = os.Stderr

		if err := gooseCmd.Run(); err != nil {
			logrus.Error("goose down migration error: ", err)
			return
		}

	},
}

var migrateFresh = &cobra.Command{
	Use:   "migrate:fresh",
	Short: "migrate fresh",
	Long:  "migrate fresh",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()

		// 1. Reset goose_db_version table first
		resetCmd := exec.Command("psql", cfg.DB_POSTGRES_DSN, "-c", "DROP TABLE IF EXISTS goose_db_version;")
		resetCmd.Stdout = os.Stdout
		resetCmd.Stderr = os.Stderr

		if err := resetCmd.Run(); err != nil {
			logrus.Error("Error resetting goose_db_version: ", err)
			// Continue anyway, as the table might not exist
		}

		// 2. Run migrations
		upCmd := exec.Command("goose", "-dir", migrationDir, "postgres", cfg.DB_POSTGRES_DSN, "up")
		upCmd.Stdout = os.Stdout
		upCmd.Stderr = os.Stderr

		if err := upCmd.Run(); err != nil {
			logrus.Error("goose up migration error: ", err)
			return
		}

		// 3. Reset seeders version table
		resetSeedCmd := exec.Command("psql", cfg.DB_POSTGRES_DSN, "-c", "DROP TABLE IF EXISTS goose_db_version_seeder;")
		resetSeedCmd.Stdout = os.Stdout
		resetSeedCmd.Stderr = os.Stderr

		if err := resetSeedCmd.Run(); err != nil {
			logrus.Error("Error resetting goose_db_version_seeder: ", err)
			// Continue anyway
		}

		// 4. Run seeders
		seedDir := "database/seeder"
		seedCmd := exec.Command("goose", "-dir", seedDir, "postgres", cfg.DB_POSTGRES_DSN, "up")
		seedCmd.Stdout = os.Stdout
		seedCmd.Stderr = os.Stderr

		if err := seedCmd.Run(); err != nil {
			logrus.Error("goose seed error: ", err)
			return
		}

		logrus.Info("Database migration and seeding completed successfully")
	},
}
