package commit

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
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
		return nil, err
	}

	return (strings.Split(response, "\n")), nil
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
	return fmt.Sprintf(`
        You are a Git commit message generator. Based on the following git diff of staged changes:

        %s

        **Criteria:**

        1. **Format:** Each commit message must follow the conventional commits format,
        which is '<type>(<scope>): <description>'.
        2. **Relevance:** Avoid mentioning a module name unless it's directly relevant
        to the change.
        3. **Clarity and Conciseness:** Each message should clearly and concisely convey
        the change made.

        **Commit Message Examples:**

        fix(app): add password regex pattern
        test(unit): add new test cases
        style: remove unused imports
        refactor(pages): extract common code to 'utils/wait.ts'

        **Recent Commits on Repo for Reference:**

        %s

        **Instructions:**

        - Take a moment to understand the changes made in the diff.
        - Think about the impact of these changes on the project (e.g., bug fixes, new
        features, performance improvements, code refactoring, documentation updates).
        - Generate 10 different commit messages that accurately describe these changes.
        - Abstract the changes to a higher level and not just describe the code changes.
        - Each message should be helpful to someone reading the project's history.

        IMPORTANT: Return only the commit messages, one per line, without any numbering,
        bullet points, or additional text. Do not include any explanations or markdown formatting.
        `,
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
