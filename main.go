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

	x := getGitBlame(args[1])

	scanner := bufio.NewScanner(strings.NewReader(x))
	for scanner.Scan() {
		line := scanner.Text()
		commitHash := getCommitHash(line)

		if cantGetPullRequestNum(commitHash) {
			fmt.Println(line)
		} else {
			gitShowOneline := getGitShowOneline(commitHash)
			if !isMergePullRequest(gitShowOneline) {
				fmt.Println(line)
			} else {
				x := strings.Replace(line, commitHash, getPullRequestNum(gitShowOneline, len(commitHash)), -1)
				fmt.Println(x)
			}

		}
	}
}
func getGitBlame(filename string) string {
	out, err := exec.Command("git", "blame", "--first-parent", filename).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	x := string(out)
	return x
}

func isMergePullRequest(gitShowOneline string) bool {
	return strings.Contains(gitShowOneline, "Merge pull request")
}

func cantGetPullRequestNum(commitHash string) bool {
	return strings.Contains(commitHash, "^")
}

func getCommitHash(line string) string {
	return strings.Split(line, " ")[0]
}

func getPullRequestNum(gitShowOneline string, commitHashlen int) string {
	x := strings.Split(gitShowOneline, " ")[4]

	return fmt.Sprintf("%"+strconv.Itoa(commitHashlen)+"s", x)
}

func getGitShowOneline(commitHash string) string {
	out, err := exec.Command("git", "show", "--oneline", commitHash).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(out)
}
