package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zored/edit/src/service/files"
	"github.com/zored/edit/src/service/navigation"
	"os"
)

var (
	file         *string
	line, column *int
	cfgFile      string
	rootCmd      = &cobra.Command{
		Use:   "reformat-tree",
		Short: "Reformat tree structure",
		Long: `Reformat tree structure on some position of file.

For example, you can turn line of parameters into column.

Or you can make one-line objects if they are small enough.
`,
		Run: func(cmd *cobra.Command, args []string) {
			config := files.NewFileFormatConfig(
				*file,
				navigation.NewPosition(*line, *column), // TODO: fix that position is +1.
			)
			if err := files.NewFileFormatter().Format(config); err != nil {
				panic(err)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Global flag:
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.edit.yaml)")

	// Local flag:
	file = rootCmd.Flags().StringP("file", "f", "", "file with tree structure")
	line = rootCmd.Flags().IntP("line", "l", 0, "file line where tree structure is")
	column = rootCmd.Flags().IntP("column", "c", 0, "file column on line where tree structure is")
}

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
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
