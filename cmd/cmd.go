package cmd

import (
	"github.com/spf13/cobra"
)

const (
	// Args
	ArgVerbose   = "verbose"
	ArgDirectory = "directory"
	ArgClass     = "class"
	ArgScript    = "script"
)

var (
	appGraph AppGraph
	app      map[string]Smali

	varClass     string
	varDirectory string
	varScript    string

	rootCmd = &cobra.Command{
		Use:   "gomli",
		Short: "A go parser for smali code",
	}

	findCmd = &cobra.Command{
		Use:   "find",
		Short: "find call chains XREFs for a particular class",

		Run: func(cmd *cobra.Command, args []string) {
			ReadDirectory(varDirectory)
			Find(varClass)
		},
	}

	replaceCmd = &cobra.Command{
		Use:   "replace",
		Short: "replace the values of const-strings globally in the application",

		Run: func(cmd *cobra.Command, args []string) {
			ReadDirectory(varDirectory)
			ReplaceArray(varScript)
		},
	}
)

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().BoolP(ArgVerbose, "v", false, "enable verbose/debug mode")

	findCmd.Flags().StringVarP(&varClass, ArgClass, "c", "", "Name of the class to search for XREF's")
	findCmd.Flags().StringVarP(&varDirectory, ArgDirectory, "d", "", "path to directory containing decompiled smali code")
	findCmd.MarkFlagRequired(ArgClass)
	findCmd.MarkFlagRequired(ArgDirectory)

	replaceCmd.Flags().StringVarP(&varScript, ArgScript, "s", "", "Javascript file that will perform the comparing and transformation")
	replaceCmd.Flags().StringVarP(&varDirectory, ArgDirectory, "d", "", "path to directory containing decompiled smali code")
	replaceCmd.MarkFlagRequired(ArgScript)
	replaceCmd.MarkFlagRequired(ArgDirectory)

	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(replaceCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
