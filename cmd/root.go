package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// Set up logger set to stdout.
func setupLoggerCtx() context.Context {
	zl := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	logger := zerologr.New(&zl)

	return logr.NewContext(context.Background(), logger.WithName("planner"))
}

func CmdContext(cmd *cobra.Command) (context.Context, logr.Logger) {
	ctx := cmd.Context()
	if cmd.Flag("verbose").Value.String() == "false" {
		ctx = logr.NewContext(cmd.Context(), logr.Discard())
	}

	log := logr.FromContextOrDiscard(ctx)

	return ctx, log
}

func New() *cobra.Command {
	name := filepath.Base(os.Args[0])

	rootCmd := &cobra.Command{
		Use:   name,
		Short: "the ultimate lightweight planner",
	}

	rootCmd.PersistentFlags().StringP("config", "c", "./config.json", "config file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	rootCmd.AddCommand(NewConfig())
	rootCmd.AddCommand(NewShow())

	rootCmd.SetContext(setupLoggerCtx())

	return rootCmd
}
