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

	out, err := exec.Command("git", "blame", "--first-parent", args[1]).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	x := string(out)

	scanner := bufio.NewScanner(strings.NewReader(x))
	for scanner.Scan() {
		line := scanner.Text()
		commitHash := getCommitHash(line)

		if cantGetPullRequestNum(commitHash) {
			fmt.Println(line)
		} else {
			isMergePullRequest, gitShowOneline := getGitShowOneline(commitHash)
			if !isMergePullRequest {
				fmt.Println(line)
			} else {
				x := strings.Replace(line, commitHash, getPullRequestNum(gitShowOneline, len(commitHash)), -1)
				fmt.Println(x)
			}

		}
	}
}

func cantGetPullRequestNum(commitHash string) bool {
	return strings.Contains(commitHash, "^")
}

func getCommitHash(line string) string {
	return strings.Split(line, " ")[0]
}

func getPullRequestNum(commitMessage string, commitHashlen int) string {
	x := strings.Split(commitMessage, " ")[4]

	return fmt.Sprintf("%"+strconv.Itoa(commitHashlen)+"s", x)
}

func getGitShowOneline(commit string) (bool, string) {
	out, err := exec.Command("git", "show", "--oneline", commit).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	x := string(out)
	isMergePullRequest := strings.Contains(x, "Merge pull request")
	return isMergePullRequest, x
}
