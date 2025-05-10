package migration

import "github.com/spf13/cobra"

var MigrateCmd = &cobra.Command{
	Use:   "migration",
	Short: "Migration commmands (create, up, down)",
}

func init() {
	MigrateCmd.AddCommand(createCmd, upCmd, downCmd)
}