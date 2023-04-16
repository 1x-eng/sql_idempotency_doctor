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
var verboseFlag bool

func init() {
	checkCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "Path to the directory that includes 'deploy' and 'revert' directories")
	checkCmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "", "Namespace for the SQL files")
	checkCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose mode")
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
			fmt.Printf("\nNot all SQL scripts are idempotent.\n")
			os.Exit(1)
		} else {
			fmt.Printf("\nAll SQL scripts are idempotent.\n")
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

	sqlScript := string(content)
	re := regexp.MustCompile(`--\s*@ddl:start\s+([\s\S]*?)\s+--\s*@ddl:end`)
	blocks := re.FindAllStringSubmatch(sqlScript, -1)

	if verboseFlag {
		fmt.Printf("\nFound %d blocks of PLPGSQL scripts that are due for assessment in '%s'.\n", len(blocks), filePath)
	}

	for _, block := range blocks {
		sqlStatements := strings.Split(block[1], ";")
		for _, sqlStatement := range sqlStatements {
			sqlStatement = strings.TrimSpace(sqlStatement)
			if len(sqlStatement) > 0 {
				isIdempotent := isIdempotent(sqlStatement, isDeploy)

				if verboseFlag {
					fmt.Printf("\n\n'%s' is idempotent? >>> %v\n\n", sqlStatement, isIdempotent)
				}

				if !isIdempotent {
					return false
				}
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

	sqlScript = strings.ToUpper(strings.TrimPrefix(sqlScript, "--@DDL"))

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
			if verboseFlag {
				fmt.Printf("\n\n'%s' matches pattern '%s'\n\n", sqlScript, pattern)
			}
			return true
		}
	}

	return false
}
