package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetChanges() ([]string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Fprintln(os.Stderr, "git bulunamadı:", err)
		os.Exit(1)
	}

	nameCmd := exec.Command("git", "diff", "--name-only", "HEAD")
	nameOut, err := nameCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("değişen dosyaları alırken hata: %w", err)
	}

	files := strings.Fields(string(nameOut))
	if len(files) == 0 {
		return nil, nil
	}

	var changes []string
	for _, file := range files {
		diffCmd := exec.Command("git", "diff", "HEAD", "--", file)
		diffOut, err := diffCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff for %s: %v\n", file, err)
			continue
		}
		changes = append(changes, fmt.Sprintf("=== Diff for %s ===\n%s\n", file, string(diffOut)))
	}

	return changes, nil
}

func CommitChanges(message string) error {
	addCmd := exec.Command("git", "add", ".")
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("git add işlemi sırasında hata: %w", err)
	}
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("commit işlemi sırasında hata: %w", err)
	}
	return nil
}
