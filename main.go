package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	if len(output) == 0 {
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
		}
		file := strings.TrimSpace(line[3:])
		diffCmd := exec.Command("git", "diff", file)
		diffOutput, err := diffCmd.Output()
		if err != nil {
			fmt.Printf("Error getting diff for %s: %v\n", file, err)
			continue
		}
		fmt.Println(string(diffOutput))
	}
}
