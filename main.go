package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("no args given")
		os.Exit(1)
	}

	cached := map[string]string{}

	gitBlameResult := getGitBlame(args[1])

	scanner := bufio.NewScanner(strings.NewReader(gitBlameResult))
	for scanner.Scan() {
		line := scanner.Text()
		commitHash := getCommitHash(line)

		if pullRequest, ok := cached[commitHash]; ok {
			replacedLine := strings.Replace(line, commitHash, pullRequest, -1)
			fmt.Println(replacedLine)
		} else {
			pullRequest := getPullRequest(commitHash)
			replacedLine := strings.Replace(line, commitHash, pullRequest, -1)
			fmt.Println(replacedLine)
			cached[commitHash] = pullRequest
		}

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
	re := regexp.MustCompile("See merge request !([0-9]+)|Merge pull request #([0-9]+)")
	res := re.FindString(getGitShow(commitHash))
	sl := strings.Split(res, " ")
	pullRequestNum := ""
	if len(sl) > 1 {
		pullRequestNum = (sl[len(sl)-1])
	} else {
		pullRequestNum = commitHash
	}
	return fmt.Sprintf("%"+strconv.Itoa(len(commitHash))+"s", pullRequestNum)
}

func getGitShow(commitHash string) string {
	out, err := exec.Command("git", "show", commitHash).Output()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return string(out)
}
