package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path"
)

var writeToDiskFlag bool

func init() {
	addCmd.PersistentFlags().BoolVarP(&writeToDiskFlag, "", "y", false, "write files to disk")
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [template name]",
	Short: "Copy files into current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("needs at least one name for the template"))
		}
		for _, name := range args {
			assets := templateMap[name]
			if assets != nil && len(assets) > 0 {
				copyAssets(assets)
				break
			}
		}
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var validTemplateNames []string
		for k := range templateMap {
			validTemplateNames = append(validTemplateNames, k)
		}
		return validTemplateNames, cobra.ShellCompDirectiveNoFileComp
	},
}

func copyAssets(assets []Asset) {
	for _, asset := range assets {
		if writeToDiskFlag == false {
			if asset.Type == "file" {
				fmt.Fprintf(os.Stdout, "Copy local file %s to %s\n", asset.Source, asset.Destination)
			} else if asset.Type == "remote" {
				fmt.Fprintf(os.Stdout, "Download remote file from %s to %s", asset.Source, asset.Destination)
			}
		} else {
			err := createDirIfNotExist(asset.Destination)
			cobra.CheckErr(err)
			if asset.Type == "file" {
				err := copyFile(asset.Source, asset.Destination)
				cobra.CheckErr(err)
			} else if asset.Type == "remote" {
				err := downloadFile(asset.Source, asset.Destination)
				cobra.CheckErr(err)
			}
		}
	}
}

func createDirIfNotExist(destination string) error {
	destinationDir := path.Dir(destination)
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destinationDir, 0754); err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(url string, dst string) error {
	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
		}
	}(destination)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)
	if _, err := io.Copy(destination, response.Body); err != nil {
		return err
	}
	return err
}

func copyFile(src string, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
		}
	}(source)
	destination, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
		}
	}(destination)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1000)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}
