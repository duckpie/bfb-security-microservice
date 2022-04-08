package main

import (
	"github.com/duckpie/bfb-security-microservice/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(cmd.NewRootCmd().Execute())
}
