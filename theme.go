package main

import "github.com/gdamore/tcell"

type Theme int

const (
	Dark Theme = iota
	Light
	SolidDark
	SolidLight
)

type TuiTheme struct {
	bg             tcell.Color
	fg             tcell.Color
	treeSelFgColor tcell.Color
	treeSelBgColor tcell.Color
	log1Bg         tcell.Color
	log1Fg         tcell.Color
}

func NewTuiTheme(theme Theme) *TuiTheme {
	switch theme {
	case Light:
		return NewLightTuiTheme()
	case SolidDark:
		return NewSolidDarkTuiTheme()
	case SolidLight:
		return NewSolidLightTuiTheme()
	default:
		return NewDarkTuiTheme()
	}
}

func NewDarkTuiTheme() *TuiTheme {
	return &TuiTheme{
		bg:             tcell.ColorBlack,
		fg:             tcell.ColorWhite,
		treeSelBgColor: tcell.ColorGreen,
		treeSelFgColor: tcell.ColorWhite,
		log1Bg:         tcell.ColorBlue,
		log1Fg:         tcell.ColorWhite,
	}
}

func NewLightTuiTheme() *TuiTheme {
	return &TuiTheme{
		bg:             tcell.ColorWhite,
		fg:             tcell.ColorBlack,
		treeSelBgColor: tcell.ColorGreen,
		treeSelFgColor: tcell.ColorWhite,
		log1Bg:         tcell.ColorBlue,
		log1Fg:         tcell.ColorWhite,
	}
}

func NewSolidDarkTuiTheme() *TuiTheme {
	return &TuiTheme{
		bg:             tcell.ColorBlack,
		fg:             tcell.ColorWhite,
		treeSelBgColor: tcell.ColorWhite,
		treeSelFgColor: tcell.ColorBlack,
		log1Bg:         tcell.ColorWhite,
		log1Fg:         tcell.ColorBlack,
	}
}

func NewSolidLightTuiTheme() *TuiTheme {
	return &TuiTheme{
		bg:             tcell.ColorWhite,
		fg:             tcell.ColorBlack,
		treeSelBgColor: tcell.ColorBlack,
		treeSelFgColor: tcell.ColorWhite,
		log1Bg:         tcell.ColorBlack,
		log1Fg:         tcell.ColorWhite,
	}
}
