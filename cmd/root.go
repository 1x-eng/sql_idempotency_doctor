/*
Copyright © 2023 Pruthvi Kumar itspruthvikumar@gmail.com

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sql_idempotency_doctor",
	Short: "A CLI tool to check the idempotency of SQL scripts in PostgreSQL",
	Long: `sql_idempotency_doctor is a CLI tool that checks the idempotency of SQL scripts in PostgreSQL.
It verifies if the PLPGSQL scripts in the specified deploy and revert directories (usually, created & maintained by sqitch) are idempotent.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
