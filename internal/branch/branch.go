package branch

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/1orzero/git-helper-cli/internal/config"
	"github.com/1orzero/git-helper-cli/internal/utils"
	"github.com/koki-develop/go-fzf"
	"github.com/tmc/langchaingo/llms"
)

func GenerateAndCleanBranchNames(llm llms.Model, description string, cfg config.Config) []string {
	branchNames, err := generateBranchNames(llm, description, cfg)
	if err != nil {
		utils.HandleError("Error generating branch names", err)
	}
	return cleanBranchNames(branchNames)
}

func generateBranchNames(llm llms.Model, description string, cfg config.Config) ([]string, error) {
	date := time.Now().Format("20060102")
	formattedDescription := formatDescription(description, cfg.Branch.DescriptionFormat, cfg.Branch.MaxDescriptionLength)

	branchPattern := strings.ReplaceAll(cfg.Branch.Pattern, "${date}", date)
	branchPattern = strings.ReplaceAll(branchPattern, "${description}", formattedDescription)

	prompt := fmt.Sprintf(`Generate %d git branch names based on this description: "%s".
	The branch names should follow this format: %s
	The description part should be concise and use hyphens instead of spaces.
	Each branch name should be unique.
	IMPORTANT: Return only the branch names, one per line, without any numbering or bullet points.`, cfg.Branch.NumSuggestions, description, branchPattern)

	ctx := context.Background()
	response, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating branch names from LLM: %v\n", err)
		return nil, err
	}

	// Remove thinking tags from the response
	cleanedResponse := utils.RemoveThinkingTags(response)

	return strings.Split(cleanedResponse, "\n"), nil
}

func formatDescription(description, format string, maxLength int) string {
	// Truncate description if it's too long
	if len(description) > maxLength {
		description = description[:maxLength]
	}

	// Convert to lowercase
	description = strings.ToLower(description)

	// Replace spaces with hyphens for kebab-case
	if format == "kebab-case" {
		description = strings.ReplaceAll(description, " ", "-")
	}

	// Remove any characters that aren't alphanumeric or hyphens
	description = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, description)

	return description
}

func cleanBranchNames(branchNames []string) []string {
	cleanedNames := make([]string, 0)
	for _, name := range branchNames {
		name = strings.TrimSpace(name)
		name = strings.TrimPrefix(name, fmt.Sprintf("%d. ", len(cleanedNames)+1))
		if name != "" {
			cleanedNames = append(cleanedNames, name)
		}
	}
	return cleanedNames
}

func SelectBranchName(branchNames []string) (string, error) {
	f, err := fzf.New(fzf.WithPrompt("Select a branch name: "))
	if err != nil {
		return "", fmt.Errorf("error initializing fzf: %w", err)
	}

	idx, err := f.Find(branchNames, func(i int) string {
		return branchNames[i]
	})
	if err != nil {
		// Check if user cancelled selection (pressed ESC)
		if errors.Is(err, fzf.ErrAbort) {
			return "", nil
		}
		return "", fmt.Errorf("error during selection: %w", err)
	}

	if len(idx) > 0 {
		selectedBranch := branchNames[idx[0]]
		return selectedBranch, nil
	}
	return "", nil
}
