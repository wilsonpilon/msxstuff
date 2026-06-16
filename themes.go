package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type ThemeID string

const (
	ThemeOneDark        ThemeID = "One Dark"
	ThemeDracula        ThemeID = "Dracula"
	ThemeMonokai        ThemeID = "Monokai"
	ThemeGitHubLight    ThemeID = "GitHub Light"
	ThemeSolarizedLight ThemeID = "Solarized Light"
	ThemeVSCodeLight    ThemeID = "VS Code Light"
	ThemeAntigravity    ThemeID = "Antigravity Dark"
	ThemeMSXClassic     ThemeID = "MSX Classic"
	ThemeMSXExpert      ThemeID = "MSX Expert"
	ThemeMSXCyber       ThemeID = "MSX Cyber"
)

var ThemeList = []string{
	string(ThemeOneDark),
	string(ThemeDracula),
	string(ThemeMonokai),
	string(ThemeGitHubLight),
	string(ThemeSolarizedLight),
	string(ThemeVSCodeLight),
	string(ThemeAntigravity),
	string(ThemeMSXClassic),
	string(ThemeMSXExpert),
	string(ThemeMSXCyber),
}

type ThemePalette struct {
	Background      color.Color
	Foreground      color.Color
	Primary         color.Color
	Button          color.Color
	InputBackground color.Color
	Focus           color.Color
	Selection       color.Color
	Hover           color.Color
	Separator       color.Color
	IsDark          bool
}

var Palettes = map[ThemeID]ThemePalette{
	ThemeOneDark: {
		Background:      color.RGBA{R: 40, G: 44, B: 52, A: 255},      // #282c34
		Foreground:      color.RGBA{R: 171, G: 178, B: 191, A: 255},   // #abb2bf
		Primary:         color.RGBA{R: 97, G: 175, B: 239, A: 255},    // #61afef
		Button:          color.RGBA{R: 44, G: 50, B: 60, A: 255},      // Escuro mas visível
		InputBackground: color.RGBA{R: 33, G: 37, B: 43, A: 255},      // #21252b
		Focus:           color.RGBA{R: 97, G: 175, B: 239, A: 128},    // Translúcido
		Selection:       color.RGBA{R: 62, G: 68, B: 81, A: 255},      // #3e4451
		Hover:           color.RGBA{R: 53, G: 59, B: 69, A: 255},      // Hover sutil
		Separator:       color.RGBA{R: 27, G: 30, B: 36, A: 255},
		IsDark:          true,
	},
	ThemeDracula: {
		Background:      color.RGBA{R: 40, G: 42, B: 54, A: 255},      // #282a36
		Foreground:      color.RGBA{R: 248, G: 248, B: 242, A: 255},   // #f8f8f2
		Primary:         color.RGBA{R: 189, G: 147, B: 249, A: 255},   // #bd93f9
		Button:          color.RGBA{R: 55, G: 58, B: 74, A: 255},      // Mais claro que bg
		InputBackground: color.RGBA{R: 30, G: 31, B: 41, A: 255},      // #1e1f29
		Focus:           color.RGBA{R: 189, G: 147, B: 249, A: 128},
		Selection:       color.RGBA{R: 68, G: 71, B: 90, A: 255},      // #44475a
		Hover:           color.RGBA{R: 68, G: 71, B: 90, A: 180},
		Separator:       color.RGBA{R: 19, G: 20, B: 26, A: 255},
		IsDark:          true,
	},
	ThemeMonokai: {
		Background:      color.RGBA{R: 39, G: 40, B: 34, A: 255},      // #272822
		Foreground:      color.RGBA{R: 248, G: 248, B: 242, A: 255},   // #f8f8f2
		Primary:         color.RGBA{R: 166, G: 226, B: 46, A: 255},    // #a6e22e (Verde clássico)
		Button:          color.RGBA{R: 62, G: 62, B: 55, A: 255},      // Cinza Monokai
		InputBackground: color.RGBA{R: 30, G: 31, B: 28, A: 255},      // #1e1f1c
		Focus:           color.RGBA{R: 249, G: 38, B: 114, A: 128},    // Rosa Monokai
		Selection:       color.RGBA{R: 73, G: 72, B: 62, A: 255},      // #49483e
		Hover:           color.RGBA{R: 73, G: 72, B: 62, A: 180},
		Separator:       color.RGBA{R: 19, G: 20, B: 17, A: 255},
		IsDark:          true,
	},
	ThemeGitHubLight: {
		Background:      color.RGBA{R: 255, G: 255, B: 255, A: 255},   // #ffffff
		Foreground:      color.RGBA{R: 36, G: 41, B: 47, A: 255},      // #24292f
		Primary:         color.RGBA{R: 9, G: 105, B: 218, A: 255},     // #0969da (Azul GitHub)
		Button:          color.RGBA{R: 246, G: 248, B: 250, A: 255},   // #f6f8fa
		InputBackground: color.RGBA{R: 243, G: 246, B: 249, A: 255},   // Cinza muito claro
		Focus:           color.RGBA{R: 9, G: 105, B: 218, A: 80},
		Selection:       color.RGBA{R: 173, G: 214, B: 255, A: 255},   // #add6ff
		Hover:           color.RGBA{R: 234, G: 238, B: 242, A: 255},
		Separator:       color.RGBA{R: 216, G: 222, B: 228, A: 255},
		IsDark:          false,
	},
	ThemeSolarizedLight: {
		Background:      color.RGBA{R: 253, G: 246, B: 227, A: 255},   // #fdf6e3
		Foreground:      color.RGBA{R: 88, G: 110, B: 117, A: 255},    // #586e75
		Primary:         color.RGBA{R: 181, G: 137, B: 0, A: 255},     // #b58900 (Solarized yellow/gold)
		Button:          color.RGBA{R: 238, G: 232, B: 213, A: 255},   // #eee8d5
		InputBackground: color.RGBA{R: 238, G: 232, B: 213, A: 255},   // #eee8d5
		Focus:           color.RGBA{R: 38, G: 139, B: 210, A: 128},    // Solarized blue focus
		Selection:       color.RGBA{R: 147, G: 161, B: 161, A: 100},   // #93a1a1
		Hover:           color.RGBA{R: 227, G: 221, B: 202, A: 255},
		Separator:       color.RGBA{R: 207, G: 200, B: 179, A: 255},
		IsDark:          false,
	},
	ThemeVSCodeLight: {
		Background:      color.RGBA{R: 243, G: 243, B: 243, A: 255},   // #f3f3f3
		Foreground:      color.RGBA{R: 51, G: 51, B: 51, A: 255},      // #333333
		Primary:         color.RGBA{R: 0, G: 122, B: 204, A: 255},     // #007acc (VS Code Blue)
		Button:          color.RGBA{R: 225, G: 223, B: 221, A: 255},   // #e1dfdd
		InputBackground: color.RGBA{R: 255, G: 255, B: 255, A: 255},   // Branco puro
		Focus:           color.RGBA{R: 0, G: 122, B: 204, A: 80},
		Selection:       color.RGBA{R: 173, G: 214, B: 255, A: 255},   // #add6ff
		Hover:           color.RGBA{R: 218, G: 218, B: 218, A: 255},
		Separator:       color.RGBA{R: 204, G: 204, B: 204, A: 255},
		IsDark:          false,
	},
	ThemeAntigravity: {
		Background:      color.RGBA{R: 15, G: 16, B: 21, A: 255},      // #0f1015
		Foreground:      color.RGBA{R: 227, G: 227, B: 227, A: 255},   // #e3e3e3
		Primary:         color.RGBA{R: 168, G: 127, B: 251, A: 255},   // #a87ffb (Gemini Violet)
		Button:          color.RGBA{R: 30, G: 30, B: 36, A: 255},      // #1e1e24
		InputBackground: color.RGBA{R: 28, G: 28, B: 34, A: 255},      // #1c1c22
		Focus:           color.RGBA{R: 168, G: 127, B: 251, A: 128},
		Selection:       color.RGBA{R: 42, G: 43, B: 54, A: 255},      // #2a2b36
		Hover:           color.RGBA{R: 34, G: 35, B: 44, A: 255},      // #22232c
		Separator:       color.RGBA{R: 45, G: 45, B: 56, A: 255},      // #2d2d38
		IsDark:          true,
	},
	ThemeMSXClassic: {
		Background:      color.RGBA{R: 26, G: 16, B: 92, A: 255},      // #1a105c (Dark Blue MSX 1)
		Foreground:      color.RGBA{R: 255, G: 255, B: 255, A: 255},   // Branco
		Primary:         color.RGBA{R: 101, G: 219, B: 239, A: 255},   // #65dbef (Ciano MSX 1)
		Button:          color.RGBA{R: 89, G: 85, B: 224, A: 255},      // #5955e0 (Medium Blue MSX 1)
		InputBackground: color.RGBA{R: 18, G: 10, B: 69, A: 255},      // Deep Blue
		Focus:           color.RGBA{R: 128, G: 118, B: 241, A: 255},   // #8076f1 (Light Blue MSX 1)
		Selection:       color.RGBA{R: 89, G: 85, B: 224, A: 255},
		Hover:           color.RGBA{R: 40, G: 33, B: 122, A: 255},
		Separator:       color.RGBA{R: 18, G: 9, B: 62, A: 255},
		IsDark:          true,
	},
	ThemeMSXExpert: {
		Background:      color.RGBA{R: 26, G: 26, B: 26, A: 255},      // #1a1a1a (Dark Charcoal)
		Foreground:      color.RGBA{R: 204, G: 204, B: 204, A: 255},   // #cccccc (Grey)
		Primary:         color.RGBA{R: 62, G: 184, B: 73, A: 255},      // #3eb849 (Medium Green MSX 1)
		Button:          color.RGBA{R: 51, G: 51, B: 51, A: 255},      // #333333
		InputBackground: color.RGBA{R: 17, G: 17, B: 17, A: 255},      // #111111
		Focus:           color.RGBA{R: 116, G: 208, B: 125, A: 128},   // #74d07d (Light Green MSX 1)
		Selection:       color.RGBA{R: 51, G: 85, B: 56, A: 255},      // Dark Green Highlight
		Hover:           color.RGBA{R: 43, G: 43, B: 43, A: 255},
		Separator:       color.RGBA{R: 34, G: 34, B: 34, A: 255},
		IsDark:          true,
	},
	ThemeMSXCyber: {
		Background:      color.RGBA{R: 19, G: 9, B: 25, A: 255},       // #130919 (Deep Magenta/Purple Black)
		Foreground:      color.RGBA{R: 101, G: 219, B: 239, A: 255},   // #65dbef (Ciano MSX 1)
		Primary:         color.RGBA{R: 183, G: 102, B: 181, A: 255},   // #b766b5 (Magenta MSX 1)
		Button:          color.RGBA{R: 50, G: 31, B: 56, A: 255},      // #321f38
		InputBackground: color.RGBA{R: 8, G: 4, B: 11, A: 255},        // Deep black-purple
		Focus:           color.RGBA{R: 101, G: 219, B: 239, A: 128},
		Selection:       color.RGBA{R: 78, G: 30, B: 85, A: 255},      // Magenta Selection
		Hover:           color.RGBA{R: 39, G: 20, B: 45, A: 255},
		Separator:       color.RGBA{R: 39, G: 20, B: 45, A: 255},
		IsDark:          true,
	},
}

type CustomFyneTheme struct {
	ID ThemeID
}

func (c *CustomFyneTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	palette, ok := Palettes[c.ID]
	if !ok {
		return theme.DefaultTheme().Color(name, variant)
	}

	// Forçar a variante do Fyne a seguir o que definimos no tema para renderizar
	// os ícones internos e elementos com alto contraste adequado.
	var targetVariant fyne.ThemeVariant = theme.VariantLight
	if palette.IsDark {
		targetVariant = theme.VariantDark
	}

	switch name {
	case theme.ColorNameBackground:
		return palette.Background
	case theme.ColorNameForeground:
		return palette.Foreground
	case theme.ColorNamePrimary:
		return palette.Primary
	case theme.ColorNameButton:
		return palette.Button
	case theme.ColorNameInputBackground:
		return palette.InputBackground
	case theme.ColorNameFocus:
		return palette.Focus
	case theme.ColorNameSeparator:
		return palette.Separator
	case theme.ColorNameSelection:
		return palette.Selection
	case theme.ColorNameHover:
		return palette.Hover
	case theme.ColorNamePlaceHolder:
		if palette.IsDark {
			return color.RGBA{R: 120, G: 120, B: 120, A: 255}
		}
		return color.RGBA{R: 150, G: 150, B: 150, A: 255}
	}

	return theme.DefaultTheme().Color(name, targetVariant)
}

func (c *CustomFyneTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (c *CustomFyneTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c *CustomFyneTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func GetTheme(id string) fyne.Theme {
	tid := ThemeID(id)
	if _, ok := Palettes[tid]; ok {
		return &CustomFyneTheme{ID: tid}
	}
	// Padrão do Fyne caso o tema seja inválido
	return theme.DefaultTheme()
}
