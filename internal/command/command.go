package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goreactnative",
	Short: "A package to fix the grade issue of apk builds.",
	Long:  "A package to fix the android.properties or gradle files in order for me to make a reactnative build locally.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
