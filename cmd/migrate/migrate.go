package migrate

import "github.com/spf13/cobra"

var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migration commmands (create, up, down)",
}

func init() {
	MigrateCmd.AddCommand(createCmd, upCmd, downCmd)
}