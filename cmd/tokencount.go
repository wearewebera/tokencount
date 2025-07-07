package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wearewebera/tokencount/tokenizer"
)

type Options struct {
	Model      string
	File       string
	Verbose    bool
	JsonOutput bool
}

type Result struct {
	Text       string                 `json:"text,omitempty"`
	TokenCount int                    `json:"token_count"`
	Model      string                 `json:"model"`
	ModelInfo  string                 `json:"model_info"`
	Details    []tokenizer.TokenInfo  `json:"details,omitempty"`
}

func Execute(text string, opts Options) error {
	// Determine the model
	var model tokenizer.Model
	switch strings.ToLower(opts.Model) {
	case "simple":
		model = tokenizer.ModelSimple
	case "gpt-3.5", "gpt3.5", "gpt35":
		model = tokenizer.ModelGPT35
	case "gpt-4", "gpt4":
		model = tokenizer.ModelGPT4
	case "claude":
		model = tokenizer.ModelClaude
	default:
		model = tokenizer.ModelGPT4
	}

	t := tokenizer.New(model)

	// Get the text to process
	var input string
	if text != "" {
		input = text
	} else if opts.File != "" {
		content, err := readFile(opts.File)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		input = content
	} else {
		// Read from stdin
		content, err := readStdin()
		if err != nil {
			return fmt.Errorf("error reading from stdin: %w", err)
		}
		input = content
	}

	if input == "" {
		return fmt.Errorf("no input text provided")
	}

	// Calculate tokens
	tokenCount := t.EstimateTokens(input)
	
	result := Result{
		TokenCount: tokenCount,
		Model:      string(model),
		ModelInfo:  t.GetModelInfo(),
	}

	if opts.Verbose {
		result.Details = t.TokenizeVerbose(input)
		if !opts.JsonOutput && len(input) <= 1000 {
			result.Text = input
		}
	}

	// Output results
	if opts.JsonOutput {
		return outputJSON(result)
	}

	return outputText(result, opts.Verbose)
}

func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func readStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", fmt.Errorf("no input provided (use -f for file or pipe content)")
	}

	reader := bufio.NewReader(os.Stdin)
	var content strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return "", err
		}

		content.Write(line)

		if err == io.EOF {
			break
		}
	}

	return content.String(), nil
}

func outputJSON(result Result) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func outputText(result Result, verbose bool) error {
	fmt.Printf("Tokens: %d\n", result.TokenCount)
	fmt.Printf("Model: %s (%s)\n", result.Model, result.ModelInfo)

	if verbose && len(result.Details) > 0 {
		fmt.Println("\nToken breakdown:")
		fmt.Println(strings.Repeat("-", 50))
		
		for _, detail := range result.Details {
			fmt.Printf("%-20s %-12s %d token(s)\n", 
				truncate(detail.Text, 20), 
				detail.Type, 
				detail.Tokens)
		}
		
		if result.Text != "" {
			fmt.Println("\nOriginal text:")
			fmt.Println(strings.Repeat("-", 50))
			fmt.Println(result.Text)
		}
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}