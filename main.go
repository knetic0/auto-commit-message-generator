package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Fprintln(os.Stderr, "git bulunamadı:", err)
		os.Exit(1)
	}

	revCmd := exec.Command("git", "rev-list", "--count", "HEAD")
	revOut, err := revCmd.Output()
	if err != nil {
		panic(fmt.Errorf("commit sayısını alırken hata: %w", err))
	}
	countStr := strings.TrimSpace(string(revOut))
	count, err := strconv.Atoi(countStr)
	if err != nil {
		panic(fmt.Errorf("commit sayısı parse edilemedi: %w", err))
	}

	var listCmd *exec.Cmd
	if count <= 1 {
		listCmd = exec.Command(
			"git", "diff-tree",
			"--no-commit-id", "--name-only", "-r", "HEAD",
		)
	} else {
		listCmd = exec.Command(
			"git", "diff", "--name-only", "HEAD^", "HEAD",
		)
	}

	listOut, err := listCmd.Output()
	if err != nil {
		panic(fmt.Errorf("değişen dosyaları listelerken hata: %w", err))
	}
	if len(listOut) == 0 {
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(listOut))
	for scanner.Scan() {
		file := strings.TrimSpace(scanner.Text())
		if file == "" {
			continue
		}

		var diffCmd *exec.Cmd
		if count <= 1 {
			diffCmd = exec.Command("git", "show", "HEAD", "--", file)
		} else {
			diffCmd = exec.Command("git", "diff", "HEAD^", "HEAD", "--", file)
		}

		diffOut, err := diffCmd.Output()
		if err != nil {
			fmt.Printf("Error getting diff for %s: %v\n", file, err)
			continue
		}

		fmt.Printf("=== Diff for %s ===\n%s\n", file, string(diffOut))
	}
}
