package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no args given")
		os.Exit(1)
	}

	gitBlameResult := getGitBlame(args[1])

	scanner := bufio.NewScanner(strings.NewReader(gitBlameResult))
	for scanner.Scan() {
		line := scanner.Text()
		commitHash := getCommitHash(line)
		pullRequest := getPullRequest(commitHash)

		replacedLine := strings.Replace(line, commitHash, pullRequest, -1)
		fmt.Println(replacedLine)

	}
}

func getGitBlame(filename string) string {
	out, err := exec.Command("git", "blame", "--first-parent", filename).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(out)
}

func getCommitHash(line string) string {
	return strings.Split(line, " ")[0]
}

func getPullRequest(commitHash string) string {
	gitShowOneline := getGitShowOneline(commitHash)
	pullRequestNum := ""
	if strings.Contains(gitShowOneline, "Merge pull request") {
		pullRequestNum = strings.Split(gitShowOneline, " ")[4]
	} else {
		pullRequestNum = commitHash
	}

	return fmt.Sprintf("%"+strconv.Itoa(len(commitHash))+"s", pullRequestNum)
}

func getGitShowOneline(commitHash string) string {
	out, err := exec.Command("git", "show", "--oneline", commitHash).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(out)
}
