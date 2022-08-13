package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
)

func NewConfig() *cobra.Command {
	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "Database configuration",
		Run:          printConfig,
		SilenceUsage: true,
	}

	return configCmd
}

// Load json file contents and unmarshal into given pointer.
func LoadJSON(filename string, v any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Load config file into object.
func LoadConfig(ctx context.Context, cfgFile string) (*DBInfo, error) {
	log := logr.FromContextOrDiscard(ctx)
	info := new(DBInfo)

	if err := LoadJSON(cfgFile, info); err != nil {
		log.Error(err, "Could not load config file.")

		return nil, err
	}

	if err := info.Validate(ctx); err != nil {
		log.Error(err, "Invalid database information.")

		return nil, err
	}

	return info, nil
}

func printConfig(cmd *cobra.Command, args []string) {
	cfgFile := cmd.Flag("config").Value.String()
	ctx, _ := CmdContext(cmd)

	cfg, err := LoadConfig(ctx, cfgFile)
	if err != nil {
		return
	}

	c := "CURRENT CONFIGURATION"
	l := "_____________________"

	// Print current db configuration.
	fmt.Printf("\n%s\n%s\n\n", c, l)
	cfg.Print()
	fmt.Println()
}
