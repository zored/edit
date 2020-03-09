package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zored/edit/src/service/formatters"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/saver"
	"github.com/zored/edit/src/service/tokens"
	"os"
)

// TODO: refactor this stuff
var (
	line, column                         *int
	file, wrapperStart, wrapperEnd, rule *string
	cfgFile                              string
	rootCmd                              = &cobra.Command{
		Use:   "reformat-tree",
		Short: "Reformat tree structure",
		Long: `Reformat tree structure on some position of file.

For example, you can turn line of parameters into column.

Or you can make one-line objects if they are small enough.
`,
		Run: func(cmd *cobra.Command, args []string) {
			options := saver.NewFileOptions(
				*file,
				navigation.NewPosition(*line, *column),
				tokens.NewWrappers(
					*wrapperStart,
					*wrapperEnd,
				),
				formatters.Rule(*rule),
			)
			err := saver.NewFileSaver().Save(options)
			handleError(err)
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	handleError(err)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("command execution error: %s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	file = rootCmd.Flags().StringP("file", "f", "", "file with tree structure")
	rule = rootCmd.Flags().StringP("rule", "r", "column", "format rule: tree, column or line")

	line = rootCmd.Flags().IntP("line", "l", 0, "file line where tree structure is")
	column = rootCmd.Flags().IntP("column", "c", 0, "file column on line where tree structure is")

	// TODO: description:
	wrapperStart = rootCmd.Flags().StringP("wrapper-start", "a", "(", "")
	wrapperEnd = rootCmd.Flags().StringP("wrapper-end", "b", ")", "")
}

// TODO: clean-up file:
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".zored_edit")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using options file:", viper.ConfigFileUsed())
	}
}
