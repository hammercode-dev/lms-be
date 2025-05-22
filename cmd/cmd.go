package cmd

import (
	"context"

	"github.com/hammer-code/lms-be/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "root cmd",
	Long:  "root cmd",
}

func Execute() {
	// set logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// load config
	config.LoadConfig()

	// Adding child commands
	rootCmd.AddCommand(createMigration, serveHttpCmd, migrateUp, migrateDown, createSeeder, seedUp)

	// cmd execute
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(context.Background(), "cmd Execute", err)
	}

}
