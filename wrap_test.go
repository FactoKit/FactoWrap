package FactoWrap_test

import (
	"testing"
	"time"

	"github.com/FactoKit/FactoWrap"
)

func TestCreateNewFactoWrap(t *testing.T) {
	factoWrap := FactoWrap.NewFactoWrap("test", []string{"test"}, "test", "test", "factorio.log")


	if _, ok := factoWrap.(FactoWrap.FactoWrap); !ok {
		t.Error("Expected factoWrap to be of type FactoWrap.FactoWrap")
	}
}

func TestFactorioWrapper(t *testing.T) {
	// RUNNING THIS REQUIRES THE FACTORIO HEADLESS SERVER BINARY. GRAB IT FROM https://factorio.com/get-download/<version>/headless/linux64
	factoWrap := FactoWrap.NewFactoWrap("./factorio/bin/x64/factorio", []string{"--start-server", "./test.zip"}, "test", "test", "factorio.log")

	go factoWrap.StartServer()
	time.Sleep(5 * time.Second)
	factoWrap.RestartServer()
	time.Sleep(5 * time.Second)
	factoWrap.SendChat("yeet haw")
	factoWrap.SendCommand("/players")
	factoWrap.SendCommand("players")
	time.Sleep(2 * time.Second)
	factoWrap.StopServer()
}