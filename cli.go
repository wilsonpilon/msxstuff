package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "msxstuffs",
	Version: fmt.Sprintf("%s (Build: %s)", version, build),
	Short:   "MSX Stuffs - Catálogo Visual de Programas MSX",
	Long: `Um catálogo visual moderno de programas e jogos de MSX da antiga série de 10 CDs da Nemesis Informática.
Este programa aceita comandos de linha de comando (CLI) via Cobra e abre a interface gráfica Fyne por padrão.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGUI()
	},
}
