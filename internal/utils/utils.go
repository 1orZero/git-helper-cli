package utils

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
)

func HandleError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
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
