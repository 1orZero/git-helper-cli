package commit

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/1orzero/git-helper-cli/internal/utils"
	"github.com/koki-develop/go-fzf"
	"github.com/tmc/langchaingo/llms"
)

// GenerateCommitMessages generates commit messages for staged changes
func GenerateCommitMessages(llm llms.Model) ([]string, error) {
	// Get staged changes
	stagedChanges, err := getStagedChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to get staged changes: %w", err)
	}

	if stagedChanges == "" {
		return nil, fmt.Errorf("no staged changes found")
	}

	// Get recent commits
	recentCommits, err := getRecentCommits(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %w", err)
	}

	// Prepare prompt for LLM
	prompt := buildPrompt(stagedChanges, recentCommits)

	// Get response from LLM
	ctx := context.Background()
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)

	if err != nil {
		fmt.Printf("Error generating commit messages from LLM: %v\n", err)
		return nil, err
	}

	// Remove thinking tags from the response
	cleanedResponse := removeThinkingTags(response)

	return strings.Split(cleanedResponse, "\n"), nil
}

// removeThinkingTags removes <think>...</think> tags from the response
func removeThinkingTags(response string) string {
	// Use regex to remove thinking tags including multi-line content
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleanedResponse := re.ReplaceAllString(response, "")

	// Trim any extra whitespace that might be left
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	return cleanedResponse
}

// getStagedChanges returns the output of git diff --cached
func getStagedChanges() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// getRecentCommits returns the recent n commits
func getRecentCommits(n int) (string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n %d", n), "--pretty=format:%h %s")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// buildPrompt constructs the prompt for the LLM
func buildPrompt(stagedChanges, recentCommits string) string {
	return fmt.Sprintf(`Generate 10 git commit messages for these staged changes:

%s

FORMAT: <type>(<scope>): <description>
Types: feat, fix, docs, style, refactor, test, chore, perf
Scope: component/module name (optional)
Description: imperative, lowercase, no period

EXAMPLES:
fix(auth): validate email format before submission
feat: add dark mode toggle to settings
refactor(api): extract request retry logic
test(user): cover edge cases for login flow

RECENT COMMITS:
%s

OUTPUT RULES:
- 10 commit messages only
- One per line
- No numbers, bullets, or explanations
- Match the style of recent commits above`,
		stagedChanges, recentCommits)
}

// extractCommitMessages extracts the commit messages from the LLM response
func extractCommitMessages(response string) string {
	// Just return the raw response - we'll clean it in cleanCommitMessages
	return response
}

func SelectCommitMessage(commitMessages []string) string {
	f, err := fzf.New(fzf.WithPrompt("Select a commit message: "))
	if err != nil {
		utils.HandleError("Error initializing fzf", err)
	}

	idx, err := f.Find(commitMessages, func(i int) string {
		return commitMessages[i]
	})
	if err != nil {
		utils.HandleError("Error during selection", err)
	}

	if len(idx) > 0 {
		selectedCommit := commitMessages[idx[0]]
		fmt.Printf("Selected commit message: %s\n", selectedCommit)
		return selectedCommit
	}
	fmt.Println("No commit message selected")
	return ""
}
