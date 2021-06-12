package FactoWrap

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type factoWrap struct {
	Config Config
	Log *log.Logger
	Mwriter io.Writer
	Pipe io.WriteCloser
	Running bool
	StopMode string
	FailCount int
}

type FactoWrap interface{
	StartServer()
	StopServer()
	RestartServer()
	SendCommand(command string)
	SendChat(message string)
}

var wg sync.WaitGroup

func NewFactoWrap(executable string, launchParameters []string, modListLocation, gameName, logLocation string) FactoWrap {
	return &factoWrap{Config: Config{
		Executable: executable,
		LaunchParameters: launchParameters,
		ModListLocation: modListLocation,
		GameName: gameName,
		LogLocation: logLocation,
	},
	Log: log.New(os.Stdout, "factowrap: ", log.Lshortfile),
	}
}



func (f *factoWrap) StartServer() {
	f.Running = false
	f.StopMode = "restart" // we want to restart by default

	// load up log from the config
	f.loadLog()

	f.bootFactorio()

	wg.Wait()
}

func (f *factoWrap) bootFactorio() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		for {
			// If it is not already running, let's start it
			if !f.Running {
				f.Running = true
				
				cmd := exec.Command(f.Config.Executable, f.Config.LaunchParameters...)
				cmd.Stderr = os.Stderr
				cmd.Stdout = f.Mwriter

				f.Pipe, err = cmd.StdinPipe()
				if err != nil {
					f.Log.Fatalf("[ERROR]: could not create standard input pipe\ndetails: %s", err.Error())
				}

				// run the command now
				err = cmd.Start()

				if err != nil {
					f.Log.Fatalf("[ERROR]: could not run the Factorio executable \ndetails: %s", err.Error())
				}

				err = cmd.Wait()
				if err != nil {
					f.Log.Fatalf("[ERROR]: %s", err.Error())
				}
			}
			if f.StopMode == "stop" {
				// break out of the for loop
				break
			} else {
				// set running to false
				f.Running = false
			}
			// wait now
			time.Sleep(2 * time.Second)
		}
	}()
}

func (f *factoWrap) StopServer() {
	// set the stop mode
	f.StopMode = "stop"
	// write the save and then stop command
	io.WriteString(f.Pipe, "/server-save\n/quit\n")

	wg.Wait()
}

func (f *factoWrap) RestartServer() {
	// set the stop mode
	f.StopMode = "restart"
	// write the save and then stop command
	io.WriteString(f.Pipe, "/server-save\n/quit\n")
}

func (f *factoWrap) SendCommand(command string) {
	// Ensure the message starts with a slash (as it is a command)
	if !strings.HasPrefix(command, "/") {
		f.Log.Printf("[WARN]: invalid use of SendCommand, must start with a '/' ")
		return
	}
	
	// write the output
	io.WriteString(f.Pipe, fmt.Sprintf("%s\n", command))
}

func (f *factoWrap) SendChat(message string) {
	// Ensure that the message does not start with a slash
	message = strings.TrimPrefix(message, "/")

	// now write the output
	io.WriteString(f.Pipe, fmt.Sprintf("%s\n", message))
}