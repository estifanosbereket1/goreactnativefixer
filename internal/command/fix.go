package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	_ "embed"
)

//go:embed gradle.properties
var gradleTemplate []byte

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Commad to generate native folders and fix file",
	Long:  "Generate native files and fix the adroid.properties file of react native project",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Println("Current working directory:", dir)
		files, err := os.ReadDir(dir)
		var isRnApp bool

		for _, val := range files {
			if val.Name() == "app.json" {
				isRnApp = true
				break
			}
			isRnApp = false
		}
		if isRnApp {
			fmt.Println("The App is react native app ")
			cmd := exec.Command("npx", "expo", "prebuild", "--platform", "android")

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmd.Run(); err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			fmt.Println("Expo prebuild completed successfully!")

			afterBuild, err := os.ReadDir(dir)

			if err != nil {
				fmt.Println(err)
				panic(err)
			}

			foundAndroid := false

			for _, val := range afterBuild {
				if val.Name() == "android" {
					foundAndroid = true
					break
				}
			}

			if foundAndroid {
				fmt.Println("We Found the android folder")
			} else {
				fmt.Println("We didnt find the android folder")
				return nil
			}

			cmdJava := exec.Command("sudo", "update-alternatives", "--config", "java")

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmdJava.Run(); err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			cmdJavac := exec.Command("sudo", "update-alternatives", "--config", "javac")

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmdJavac.Run(); err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			androidDirErr := os.Chdir("android")

			if androidDirErr != nil {
				panic(androidDirErr)
			}

			androidDir, errAndroidDir := os.Getwd()
			if errAndroidDir != nil {
				panic(errAndroidDir)
			}

			androidDirFiles, androidDirFilesErr := os.ReadDir(androidDir)

			if androidDirFilesErr != nil {
				panic(androidDirFilesErr)
			}
			hasGradleFile := false
			for _, androidF := range androidDirFiles {
				fmt.Println(androidF.Name())
				if androidF.Name() == "gradle.properties" {
					hasGradleFile = true
					break
				}
			}
			if hasGradleFile {

				fullPath := filepath.Join(androidDir, "gradle.properties")
				err := os.WriteFile(fullPath, gradleTemplate, 0644)
				if err != nil {
					panic(err)
				}
				fmt.Println("Writing the gradle file completed")
			}

		} else {
			fmt.Println("Not a react native app")
		}

		return nil
	},
}

var pwd string

func init() {
	rootCmd.AddCommand(fixCmd)
}
