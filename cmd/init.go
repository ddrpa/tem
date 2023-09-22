package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path"
	"tem/boilerplate"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

// 负责将模板文件复制到当前目录
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial or refresh tem configuration and default boilerplate files",
	Long:  "Run after fresh installation or update, ensure tem directory and default boilerplate files exist",
	Run: func(cmd *cobra.Command, args []string) {
		err := createTemHomeDirIfNotExist()
		cobra.CheckErr(err)
		err = cleanUpDefaultBoilerplateDir()
		cobra.CheckErr(err)
		err = createCustomDirIfNotExist()
		cobra.CheckErr(err)
		fmt.Fprintf(os.Stdout, "Initial complete, tem home directory is %s, use 'tem completion <sh>' for shell completion\n", temHomeDir)
	},
}

func createTemHomeDirIfNotExist() error {
	if _, err := os.Stat(temHomeDir); os.IsNotExist(err) {
		if err := os.Mkdir(temHomeDir, 0754); err != nil {
			return err
		}
	}
	return nil
}

func cleanUpDefaultBoilerplateDir() error {
	if _, err := os.Stat(defaultBoilerplateDir); err == nil {
		if err := os.RemoveAll(defaultBoilerplateDir); err != nil {
			return err
		}
	}
	if err := os.Mkdir(defaultBoilerplateDir, 0754); err != nil {
		return err
	}
	entries, _ := boilerplate.DefaultBoilerplateFs.ReadDir(".")
	for _, e := range entries {
		if err := walkAndHandleEmbedResources("", e); err != nil {
			return err
		}
	}
	return nil
}

func walkAndHandleEmbedResources(parentDir string, entry fs.DirEntry) error {
	var entryPath string
	if parentDir == "" {
		entryPath = entry.Name()
	} else {
		entryPath = path.Join(parentDir, entry.Name())
	}
	if entry.IsDir() {
		if err := os.Mkdir(path.Join(defaultBoilerplateDir, entryPath), 0754); err != nil {
			return err
		}
		if entries, err := boilerplate.DefaultBoilerplateFs.ReadDir(entryPath); err == nil {
			for _, e := range entries {
				if err := walkAndHandleEmbedResources(entryPath, e); err != nil {
					return err
				}
			}
		}
	} else {
		if bytes, err := boilerplate.DefaultBoilerplateFs.ReadFile(entryPath); err == nil {
			if err := os.WriteFile(path.Join(defaultBoilerplateDir, entryPath), bytes, 0654); err != nil {
				return err
			}
		}
	}
	return nil
}

func createCustomDirIfNotExist() error {
	if _, err := os.Stat(customBoilerplateDir); os.IsNotExist(err) {
		if err := os.Mkdir(customBoilerplateDir, 0754); err != nil {
			return err
		}
	}
	if _, err := os.Stat(path.Join(customBoilerplateDir, "custom.toml")); os.IsNotExist(err) {
		if err := os.WriteFile(path.Join(customBoilerplateDir, "custom.toml"), boilerplate.CustomConfigExample, 0654); err != nil {
			return err
		}
	}
	return nil
}
