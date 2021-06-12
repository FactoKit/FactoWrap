package FactoWrap

import (
	"io"
	"os"
)

func (f *factoWrap) loadLog() {
	if err := os.Remove(f.Config.LogLocation); err != nil {
		f.Log.Printf("[WARN]: %s doesn't exist, continuing anyway\n", f.Config.LogLocation)
	}

	logging, err := os.OpenFile("factorio.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		f.Log.Fatalf("[ERROR]: could not open the log file at %s\nerror information: %s", f.Config.LogLocation, err.Error())
	}


	f.Mwriter = io.MultiWriter(logging, os.Stdout)
}