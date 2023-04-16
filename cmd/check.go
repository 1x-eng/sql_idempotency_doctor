package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var pathFlag string
var namespaceFlag string

func init() {
	checkCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to the directory that includes 'deploy' and 'revert' directories")
	checkCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace for the SQL files")
	checkCmd.MarkFlagRequired("path")
	checkCmd.MarkFlagRequired("namespace")
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if SQL scripts in deploy and revert directories are idempotent",
	Long: `This command checks if all the SQL scripts in the <namespace>.sql file in the 'deploy' and 'revert' directories
are idempotent. If even one of them is not, it returns false.`,
	Run: func(cmd *cobra.Command, args []string) {
		isIdempotent := checkSQLFiles(pathFlag, namespaceFlag)
		if !isIdempotent {
			fmt.Println("Not all SQL scripts are idempotent.")
			os.Exit(1)
		} else {
			fmt.Println("All SQL scripts are idempotent.")
		}
	},
}

func checkSQLFiles(basePath, namespace string) bool {
	deployPath := filepath.Join(basePath, "deploy", fmt.Sprintf("%s.sql", namespace))
	revertPath := filepath.Join(basePath, "revert", fmt.Sprintf("%s.sql", namespace))

	return checkSQLFile(deployPath, true) && checkSQLFile(revertPath, false)
}

func checkSQLFile(filePath string, isDeploy bool) bool {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	sqlStatements := strings.Split(string(content), ";")
	for _, sqlStatement := range sqlStatements {
		sqlStatement = strings.TrimSpace(sqlStatement)
		if len(sqlStatement) > 0 {
			isIdempotent := isIdempotent(sqlStatement, isDeploy)
			if !isIdempotent {
				return false
			}
		}
	}

	return true
}

func isIdempotent(sqlScript string, isDeploy bool) bool {
	if !isDeploy {
		// For revert scripts, just return true for now.
		// Add any specific rules if necessary.
		return true
	}

	sqlScript = strings.ToLower(sqlScript)

	idempotentPatterns := []string{
		`CREATE\s+TABLE\s+IF\s+NOT\s+EXISTS`,
		`CREATE\s+(?:UNIQUE\s+)?INDEX\s+IF\s+NOT\s+EXISTS`,
		`ALTER\s+TABLE\s+\w+\s+ADD\s+COLUMN\s+IF\s+NOT\s+EXISTS`,
		`CREATE\s+SEQUENCE\s+IF\s+NOT\s+EXISTS`,
		`CREATE\s+OR\s+REPLACE\s+(?:VIEW|FUNCTION|PROCEDURE|TRIGGER|AGGREGATE|OPERATOR|RULE|POLICY|EVENT\s+TRIGGER|LANGUAGE|EXTENSION)`,
		`CREATE\s+(?:ROLE|USER|SCHEMA|DOMAIN|CAST|COLLATION|CONVERSION|TYPE|SERVER|FOREIGN\s+TABLE|MATERIALIZED\s+VIEW|PUBLICATION|SUBSCRIPTION)\s+IF\s+NOT\s+EXISTS`,
		`CREATE\s+TEXT\s+SEARCH\s+(?:DICTIONARY|CONFIGURATION|PARSER|TEMPLATE)\s+IF\s+NOT\s+EXISTS`,
	}

	for _, pattern := range idempotentPatterns {
		regex := regexp.MustCompile(pattern)
		if regex.MatchString(sqlScript) {
			return true
		}
	}

	return false
}
