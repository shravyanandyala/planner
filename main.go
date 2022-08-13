package main

import (
	"fmt"
	"os"

	"github.com/shravyanandyala/planner/cmd"
)

func main() {
	if e := execRootCmd(); e != nil {
		os.Exit(1)
	}
}

func execRootCmd() error {
	rootCmd := cmd.New()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
