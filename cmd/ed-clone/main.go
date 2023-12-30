package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
  "time"
  "strings")

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
		cmd := exec.Command("rlwrap", "-a", os.Args[0], "--with-rlwrap")
    fmt.Println("\033[1;33m[INFO]\033[0m -", cmd)
    fmt.Println("\033[1;33m[INFO]\033[0m - Running command with rlwrap.")
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
  fileContents := "Branch: " + newBranchName + " - " + time.Now().Format("2006-01-02 15:04:05")

  // Now append this as a new line to the "branches" file
  f, err := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0644)
  f.Write([]byte(fileContents))

  if err != nil {
    fmt.Printf("Error creating new file: %v\n", err)
  }
  err = executeCommand("git", "add", ".")
  if err != nil {
    fmt.Printf("Error adding new file: %v\n", err)
  }
  
  err = executeCommand("git", "add", ".")
  err = executeCommand("git", "commit", "-m", "Added new branch: " + newBranchName)
  err = executeCommand("git", "push", "--set-upstream", "origin", newBranchName)
  err = executeCommand("git", "push")
  if err != nil {
    fmt.Printf("Error pushing new branch: %v\n", err)
  }
  // This will push this branch & the new readme file to the remote repository

	// Open the specified code editor
	switch codeEditor {
	case "vscode":
		err = executeCommand("code", ".")
	case "nvim":
		err = executeCommand("nvim", ".")
  case "lvim":
    err = executeCommand("lvim", "./")
  case "vim":
    err = executeCommand("vim", ".")
  case "kate":
    err = executeCommand("kate", ".")
	default:
		fmt.Println("Unsupported code editor")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error opening code editor: %v\n", err)
		os.Exit(1)
	}
}
