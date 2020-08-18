package main

import (
	"github.com/d1y/macapp"
)

func main() {
	macapp.Create(macapp.AppConfig{
		AppName: `xiaoya`,
		AppPath: `/Users/kozo4/cat`,
	})
	// Xiaoya.SetBinFile()
	// Xiaoya.SetIcon()
}
