package main

import (
	"fmt"
	"os"
)

var (
	version = "3.1.1"
	build   = "dev"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao executar o comando: %v\n", err)
		os.Exit(1)
	}
}
