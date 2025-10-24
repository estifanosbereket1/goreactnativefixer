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

func isReactnativeApp(files *[]os.DirEntry) bool {
	var isNative bool
	for _, val := range *files {
		if val.Name() == "app.json" {
			isNative = true
			break
		}
	}
	return isNative
}
func runPreBuild(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func findAndroidFolder(afterBuild *[]os.DirEntry) bool {
	var foundAndroid bool
	for _, val := range *afterBuild {
		if val.Name() == "android" {
			foundAndroid = true
			break
		}
	}
	return foundAndroid
}

func updateJavaVerions(cmdJava *exec.Cmd) error {
	cmdJava.Stdout = os.Stdout
	cmdJava.Stderr = os.Stderr
	cmdJava.Stdin = os.Stdin

	if err := cmdJava.Run(); err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return nil
}

func updateJavacVerions(cmdJavaC *exec.Cmd) error {
	cmdJavaC.Stdout = os.Stdout
	cmdJavaC.Stderr = os.Stderr
	cmdJavaC.Stdin = os.Stdin

	if err := cmdJavaC.Run(); err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return nil
}

func changeAndGetDir(dirName string) ([]os.DirEntry, string, error) {
	androidDirErr := os.Chdir(dirName)
	if androidDirErr != nil {
		return nil, "", androidDirErr
	}
	androidDir, errAndroidDir := os.Getwd()
	if errAndroidDir != nil {
		return nil, "", errAndroidDir
	}
	androidDirFiles, androidDirFilesErr := os.ReadDir(androidDir)

	if androidDirFilesErr != nil {
		return nil, "", androidDirFilesErr
	}
	return androidDirFiles, androidDir, nil
}

func hasGradleFile(androidDirFiles *[]os.DirEntry, fileName string) bool {
	var hasGradleFile bool

	for _, androidF := range *androidDirFiles {
		if androidF.Name() == fileName {
			hasGradleFile = true
			break
		}
	}
	return hasGradleFile
}

func writeGradleFile(fullPath string) error {
	err := os.WriteFile(fullPath, gradleTemplate, 0644)
	if err != nil {
		return err
	}
	return nil
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Commad to generate native folders and fix file",
	Long:  "Generate native files and fix the adroid.properties file of react native project",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Println("Current working directory:", dir)
		files, errFile := os.ReadDir(dir)

		if errFile != nil {
			return errFile
		}
		isRnApp := isReactnativeApp(&files)

		if isRnApp {
			fmt.Println("The App is react native app ")
			cmd := exec.Command("npx", "expo", "prebuild", "--platform", "android")

			preBildErr := runPreBuild(cmd)

			if preBildErr != nil {
				return preBildErr
			}

			fmt.Println("Expo prebuild completed successfully!")

			afterBuild, afterBuildErr := os.ReadDir(dir)

			if afterBuildErr != nil {
				fmt.Println(afterBuildErr)
				return afterBuildErr
			}

			foundAndroid := findAndroidFolder(&afterBuild)

			if foundAndroid {
				fmt.Println("We Found the android folder")
			} else {
				fmt.Println("We didnt find the android folder")
				return fmt.Errorf("couldn't get android folder")
			}

			cmdJava := exec.Command("sudo", "update-alternatives", "--config", "java")

			cmdJavaError := updateJavaVerions(cmdJava)

			if cmdJavaError != nil {
				fmt.Println(cmdJavaError)
				panic(cmdJavaError)
			}

			cmdJavac := exec.Command("sudo", "update-alternatives", "--config", "javac")

			cmdJavaCErr := updateJavacVerions(cmdJavac)

			if cmdJavaCErr != nil {
				fmt.Println(cmdJavaCErr)
				panic(cmdJavaCErr)
			}

			androidDirFiles, androidDir, getAndroidDir := changeAndGetDir("android")

			if getAndroidDir != nil {
				fmt.Println(getAndroidDir)
				return getAndroidDir
			}

			hasGradleFile := hasGradleFile(&androidDirFiles, "gradle.properties")
			if hasGradleFile {
				fullPath := filepath.Join(androidDir, "gradle.properties")
				err := writeGradleFile(fullPath)
				if err != nil {
					return err
				}
				fmt.Println("Gradle File written succesfully")
			}
		} else {
			fmt.Println("Not a react native app")
			return fmt.Errorf("this is not a react native project")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fixCmd)
}
