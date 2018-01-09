package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no args given")
		os.Exit(1)
	}

	// cached := map[string]string{}

	gitBlameResult := getGitBlame(args[1])

	scanner := bufio.NewScanner(strings.NewReader(gitBlameResult))
	for scanner.Scan() {
		line := scanner.Text()
		commitHash := getCommitHash(line)

		fmt.Println(getPullRequest(commitHash))
	}
}

func getPullRequest(commitHash string) string {
	re := regexp.MustCompile("See merge request !([0-9]+)|Merge pull request #([0-9]+)")
	res := re.FindString(getGitShowOneline(commitHash))
	sl := strings.Split(res, " ")
	if len(sl) > 1 {
		return (sl[len(sl)-1])
	}
	return commitHash
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

func getGitShowOneline(commitHash string) string {
	out, err := exec.Command("git", "show", commitHash).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(out)
}
