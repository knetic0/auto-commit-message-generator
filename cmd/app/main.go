package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/knetic0/auto-commit-message-generator/config"
	"github.com/knetic0/auto-commit-message-generator/internal/gemini"
	"github.com/knetic0/auto-commit-message-generator/internal/git"
	"golang.org/x/term"
)

func main() {
	if err := config.Init(); err != nil {
		panic(err)
	}

	geminiClient := gemini.NewGeminiClient()
	changes, err := git.GetChanges()
	if err != nil {
		panic(err)
	}
	if len(changes) == 0 {
		println("No changes detected.")
		return
	}
	changesText := "Changes:\n" + strings.Join(changes, "\n")
	request := gemini.NewGeminiRequest(changesText)
	response, err := geminiClient.GenerateCommitMessage(request)
	if err != nil {
		panic(err)
	}
	if response == nil || len(response.Candidates) == 0 {
		println("No commit message generated.")
		return
	}
	message := response.Candidates[0].Content.Parts[0].Text

	fmt.Println("Generated commit message:", message)
	fmt.Println("\nPress ENTER to commit, or ESC to cancel.")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)

	buf := make([]byte, 1)
	if _, err := os.Stdin.Read(buf); err != nil {
		panic(err)
	}

	switch buf[0] {
	case 13, 10:
		if err := git.CommitChanges(message); err != nil {
			panic(err)
		}
		fmt.Println("\n✅ Changes committed.")
	case 27:
		fmt.Println("\n❌ Commit canceled.")
	default:
		fmt.Printf("\n⚠️  Unrecognized key (code=%d), aborting.\n", buf[0])
	}
}
