package tokenizer

import (
	"regexp"
	"strings"
	"unicode"
)

type Model string

const (
	ModelSimple Model = "simple"
	ModelGPT35  Model = "gpt-3.5"
	ModelGPT4   Model = "gpt-4"
	ModelClaude Model = "claude"
)

type Tokenizer struct {
	model Model
}

func New(model Model) *Tokenizer {
	return &Tokenizer{model: model}
}

func (t *Tokenizer) EstimateTokens(text string) int {
	switch t.model {
	case ModelSimple:
		return t.simpleEstimate(text)
	case ModelGPT35, ModelGPT4:
		return t.gptEstimate(text)
	case ModelClaude:
		return t.claudeEstimate(text)
	default:
		return t.simpleEstimate(text)
	}
}

func (t *Tokenizer) simpleEstimate(text string) int {
	if text == "" {
		return 0
	}
	return (len(text) + 3) / 4
}

func (t *Tokenizer) gptEstimate(text string) int {
	if text == "" {
		return 0
	}

	tokens := 0
	
	// Split text into words and special characters
	// This regex splits on word boundaries while preserving punctuation
	parts := regexp.MustCompile(`\w+|[^\w\s]+|\s+`).FindAllString(text, -1)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		// Skip pure whitespace
		if regexp.MustCompile(`^\s+$`).MatchString(part) {
			continue
		}
		
		// Check if it's Unicode text (non-ASCII)
		hasUnicode := false
		for _, r := range part {
			if r > 127 {
				hasUnicode = true
				break
			}
		}
		
		// Unicode text (Chinese, Japanese, etc.) - each character is often a token
		if hasUnicode {
			tokens += len([]rune(part))
			continue
		}
		
		// Punctuation and special characters (each counts as 1 token)
		if regexp.MustCompile(`^[^\w\s]+$`).MatchString(part) {
			tokens++
			continue
		}
		
		// Numbers
		if regexp.MustCompile(`^\d+$`).MatchString(part) {
			if len(part) <= 3 {
				tokens++
			} else {
				tokens += (len(part) + 2) / 3
			}
			continue
		}
		
		// Regular words
		wordLen := len(part)
		if wordLen <= 5 {
			tokens++
		} else if wordLen <= 10 {
			// Medium words often split into 2 tokens
			tokens += (wordLen + 4) / 5
		} else {
			// Longer words are split into subwords
			tokens += (wordLen + 3) / 4
		}
	}
	
	return tokens
}

func (t *Tokenizer) claudeEstimate(text string) int {
	if text == "" {
		return 0
	}
	
	// Claude's tokenization is similar to GPT but with some differences
	tokens := 0
	
	// Split by whitespace and punctuation
	parts := regexp.MustCompile(`(\s+|\W+|\w+)`).FindAllString(text, -1)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		// Whitespace handling
		if strings.TrimSpace(part) == "" {
			if len(part) > 1 {
				tokens++
			}
			continue
		}
		
		// Punctuation and special characters
		if regexp.MustCompile(`^\W+$`).MatchString(part) {
			tokens += len(regexp.MustCompile(`\W`).FindAllString(part, -1))
			continue
		}
		
		// Regular words
		wordLen := len(part)
		if wordLen <= 5 {
			tokens++
		} else {
			// Estimate based on syllables/subwords
			tokens += (wordLen + 4) / 5
		}
	}
	
	return tokens
}

func (t *Tokenizer) GetModelInfo() string {
	switch t.model {
	case ModelSimple:
		return "Simple estimation (4 characters = 1 token)"
	case ModelGPT35:
		return "GPT-3.5 estimation algorithm"
	case ModelGPT4:
		return "GPT-4 estimation algorithm"
	case ModelClaude:
		return "Claude estimation algorithm"
	default:
		return "Unknown model"
	}
}

func hasUpperCase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func countUpperCase(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsUpper(r) {
			count++
		}
	}
	return count
}

// TokenizeVerbose returns a detailed breakdown of tokenization
type TokenInfo struct {
	Text   string
	Tokens int
	Type   string
}

func (t *Tokenizer) TokenizeVerbose(text string) []TokenInfo {
	var result []TokenInfo
	
	if t.model == ModelSimple {
		return []TokenInfo{{Text: text, Tokens: t.simpleEstimate(text), Type: "simple"}}
	}
	
	parts := regexp.MustCompile(`\w+|[^\w\s]+|\s+`).FindAllString(text, -1)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		info := TokenInfo{Text: part}
		
		if strings.TrimSpace(part) == "" {
			info.Type = "whitespace"
			info.Tokens = 0
		} else {
			// Check for Unicode
			hasUnicode := false
			for _, r := range part {
				if r > 127 {
					hasUnicode = true
					break
				}
			}
			
			if hasUnicode {
				info.Type = "unicode"
				info.Tokens = len([]rune(part))
			} else if regexp.MustCompile(`^[^\w\s]+$`).MatchString(part) {
				info.Type = "punctuation"
				info.Tokens = 1
			} else if regexp.MustCompile(`^\d+$`).MatchString(part) {
				info.Type = "number"
				info.Tokens = (len(part) + 2) / 3
			} else {
				info.Type = "word"
				wordLen := len(part)
				if wordLen <= 5 {
					info.Tokens = 1
				} else if wordLen <= 10 {
					info.Tokens = (wordLen + 4) / 5
				} else {
					info.Tokens = (wordLen + 3) / 4
				}
			}
		}
		
		if info.Tokens > 0 {
			result = append(result, info)
		}
	}
	
	return result
}