package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Define a flag for folder path
	folderPtr := flag.String("folder", "", "Path to the folder to watch")

	// Define a string slice to store command arguments
	var command []string

	// Define a flag to handle command arguments
	flag.Var(&command, "command", "Command to run on changes (can be repeated for multiple arguments)")
	flag.Parse()

	// Check if both folder path and command arguments are provided
	if *folderPtr == "" || len(command) == 0 {
		fmt.Println("Usage: watch -folder <folder_path> -command <command_argument> [<command_argument>...]")
		flag.PrintDefaults()
		return
	}

	// Get folder path
	folder := *folderPtr

	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Function to run the command
	runCommand := func() {
		// Concatenate command arguments
		cmdStr := fmt.Sprintf("%s", command[0]) // Start with the first argument
		for _, arg := range command[1:] {
			cmdStr += fmt.Sprintf(" %s", arg)
		}

		// Replace placeholders with actual arguments (optional)
		// You can modify the command string (cmdStr) further if needed

		// Execute the command
		cmd := exec.Command(cmdStr, fmt.Sprintf("%s", folder))
		err := cmd.Run()
		if err != nil {
			log.Printf("Error running command: %v\n", err)
		} else {
			log.Println("Command executed successfully!")
		}
	}

	// Add the folder to watch
	if err := watcher.Add(folder); err != nil {
		log.Fatal(err)
	}

	log.Printf("Watching folder: %s\n", folder)

	// Infinite loop to listen for events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op == fsnotify.Create || event.Op == fsnotify.Write || event.Op == fsnotify.Remove {
				log.Printf("Event: %v - Path: %v\n", event.Op, event.Name)
				go runCommand() // Run the command in a goroutine
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
