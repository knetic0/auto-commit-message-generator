package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetChanges() ([]string, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, fmt.Errorf("git bulunamadı: %w", err)
	}

	revParse := exec.Command("git", "rev-parse", "--show-toplevel")
	topLevelBytes, err := revParse.Output()
	if err != nil {
		return nil, fmt.Errorf("repo kök dizini bulunamadı: %w", err)
	}
	root := strings.TrimSpace(string(topLevelBytes))

	if err := os.Chdir(root); err != nil {
		return nil, fmt.Errorf("köke geçiş başarısız: %w", err)
	}

	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOut, err := statusCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git status çalıştırılamadı: %w", err)
	}
	if len(statusOut) == 0 {
		return nil, nil
	}

	scanner := bufio.NewScanner(bytes.NewReader(statusOut))
	var files []string
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
		}
		file := strings.TrimSpace(line[3:])
		file = filepath.ToSlash(file)
		files = append(files, file)
	}
	if len(files) == 0 {
		return nil, nil
	}

	var changes []string
	for _, file := range files {
		diffCmd := exec.Command("git", "diff", "HEAD", "--", file)
		diffOut, err := diffCmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "diff okunamadı %s: %v\n", file, err)
			continue
		}
		changes = append(changes, fmt.Sprintf("=== Diff for %s ===\n%s", file, string(diffOut)))
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
