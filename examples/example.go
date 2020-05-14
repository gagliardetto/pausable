package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gagliardetto/keyable"
	"github.com/gagliardetto/pausable"
)

func main() {
	pauser := pausable.New()

	// Print some indications on usage:
	fmt.Println("Press P to pause/play; Press S to see stats; press ESC or CTRL+C to quit; press CTRL+F to force-quit")

	// Create a new Keyable object and add callbacks
	// for the various characters and key combinations:
	kb := keyable.New()

	// onStat will be the callback for information for your program
	// (progress, etc.)
	onStat := func() {
		fmt.Println("some stats here")
	}
	// onQuit is a callback for when you quit the program:
	onQuit := func() {
		fmt.Println("Exiting...")
		// TODO: clean up, send exit signals, etc.
		os.Exit(0)
	}

	// Add the callbacks:
	kb.
		OnChar(
			onStat,
			'i',
		).
		OnChar(
			func() {
				paused := pauser.Toggle()
				if paused {
					fmt.Println("pause")
				} else {
					fmt.Println("play")
				}
			},
			'p',
		).
		OnKey(
			onQuit,
			keyboard.KeyEsc, keyboard.KeyCtrlC,
		).
		OnKey(
			func() {
				fmt.Println("FORCE-QUITTING...")
				os.Exit(1)
			},
			keyboard.KeyCtrlF,
		)
	// Start the keyboard listener:
	err := kb.Start()
	if err != nil {
		panic(err)
	}
	// Don't forget to stop it!
	defer kb.Stop()

	// This is some multi-task job you're running:
	for {
		// Wait for green-light (it will block if the pausable is paused, until it is un-paused):
		pauser.Wait()

		// Execute your thing:
		RunTask()
	}
}

func RunTask() error {
	fmt.Println("Task running at", time.Now().Format("15:04:05"))
	time.Sleep(time.Second)
	return nil
}
