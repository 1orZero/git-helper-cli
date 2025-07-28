package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
)

func HandleError(message string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}

func CopyToClipboard(text string) {
	if text == "" {
		return
	}
	err := clipboard.WriteAll(text)
	if err != nil {
		fmt.Printf("Error copying to clipboard: %v\n", err)
	} else {
		fmt.Println("copied to clipboard!")
	}
}

func Log(message string) {
	fmt.Println(message)
}

// RemoveThinkingTags removes <think>...</think> tags from the response
func RemoveThinkingTags(response string) string {
	// Use regex to remove thinking tags including multi-line content
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleanedResponse := re.ReplaceAllString(response, "")
	
	// Trim any extra whitespace that might be left
	cleanedResponse = strings.TrimSpace(cleanedResponse)
	
	return cleanedResponse
}
