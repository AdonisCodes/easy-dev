package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
  "strings"
)

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func isRunningWithRlwrap() bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, "rlwrap") {
			return true
		}
	}
	return false
}

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Prompt for input
  if !isRunningWithRlwrap() {
		// Re-execute the program with rlwrap
		cmd := exec.Command("rlwrap", os.Args[0], "--with-rlwrap")
    fmt.Println("[INFO] -", cmd)
    fmt.Println("[INFO] - Running command with rlwrap.")
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("Error re-executing with rlwrap:", err)
			os.Exit(1)
		}
		return
	}

	folderName := promptForInput("Enter folder name: ")
	repositoryURL := promptForInput("Enter Github url: ")
	newBranchName := promptForInput("Enter new branch name: ")
	codeEditor := promptForInput("Enter code editor (e.g., vscode): ")

	// Clone repository
	err := executeCommand("git", "clone", "https://github.com/"+repositoryURL, folderName)
	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
		os.Exit(1)
	}

	// Change into the newly created folder
	err = os.Chdir(folderName)
	if err != nil {
		fmt.Printf("Error changing into the repository directory: %v\n", err)
		os.Exit(1)
	}

	// Create and switch to a new branch
	err = executeCommand("git", "checkout", "-b", newBranchName)
	if err != nil {
		fmt.Printf("Error creating/switching to branch: %v\n", err)
		os.Exit(1)
	}

	// Open the specified code editor
	switch codeEditor {
	case "vscode":
		err = executeCommand("code", ".")
	case "nvim":
		err = executeCommand("nvim", ".")
	default:
		fmt.Println("Unsupported code editor")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error opening code editor: %v\n", err)
		os.Exit(1)
	}
}
