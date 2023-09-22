package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

type RawTemplate struct {
	Key    string     `mapstructure:"key"`
	Alias  []string   `mapstructure:"alias"`
	Assets [][]string `mapstructure:"assets"`
}

type Template struct {
	Key    string
	Assets []Asset
}

type Asset struct {
	Type        string
	Source      string
	Destination string
}

const (
	Version = "0.0.1"
)

var temHomeDir string
var defaultBoilerplateDir string
var customBoilerplateDir string
var templateMap = make(map[string][]Asset)

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	temHomeDir = path.Join(home, ".tem")
	defaultBoilerplateDir = path.Join(temHomeDir, "default")
	customBoilerplateDir = path.Join(temHomeDir, "custom")
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:     "tem",
	Short:   "Tem is a boilerplate injector",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	temHomeDir = path.Join(home, ".tem")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// load default config
	defaultConfigViper := viper.New()
	defaultConfigViper.AddConfigPath(defaultBoilerplateDir)
	defaultConfigViper.SetConfigType("toml")
	defaultConfigViper.SetConfigName("default")
	if err := defaultConfigViper.ReadInConfig(); err == nil {
		var rawTemplates []RawTemplate
		err := defaultConfigViper.UnmarshalKey("template", &rawTemplates)
		cobra.CheckErr(err)
		loadTemplate(templateMap, rawTemplates)
	} else {
		fmt.Fprintln(os.Stderr, "Parse default configuration failed, 'tem init' may fix it")
	}

	// load custom config
	customConfigViper := viper.New()
	customConfigViper.AddConfigPath(customBoilerplateDir)
	customConfigViper.SetConfigType("toml")
	customConfigViper.SetConfigName("custom")
	if err := customConfigViper.ReadInConfig(); err == nil {
		var rawTemplates []RawTemplate
		err := customConfigViper.UnmarshalKey("template", &rawTemplates)
		cobra.CheckErr(err)
		loadTemplate(templateMap, rawTemplates)
	} else {
		fmt.Fprintln(os.Stderr, "Parse customized configuration failed")
	}
}

// loadTemplate load template from raw template in config file
func loadTemplate(tMap map[string][]Asset, raw []RawTemplate) {
	for _, t := range raw {
		// skip empty key
		if t.Key == "" {
			continue
		}
		assets := make([]Asset, 0, 2)
		for _, assetTuple := range t.Assets {
			if len(assetTuple) < 2 {
				// skip invalid asset
				continue
			}
			var fileSet Asset
			if strings.HasPrefix(assetTuple[0], "https://") || strings.HasPrefix(assetTuple[0], "http://") {
				// if source start with http schema, treat as remote file
				fileSet = Asset{
					Type:        "remote",
					Source:      assetTuple[0],
					Destination: assetTuple[1],
				}
			} else {
				fileSet = Asset{
					Type:        "file",
					Source:      path.Join(temHomeDir, assetTuple[0]),
					Destination: assetTuple[1],
				}
			}
			assets = append(assets, fileSet)
		}
		tMap[t.Key] = assets
		if len(t.Alias) > 0 {
			// add alias
			for _, a := range t.Alias {
				tMap[a] = assets
			}
		}
	}
}
