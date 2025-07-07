package tokenizer

import (
	"testing"
)

func TestSimpleEstimate(t *testing.T) {
	tokenizer := New(ModelSimple)
	
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Four characters", "test", 1},
		{"Eight characters", "testtest", 2},
		{"With spaces", "hello world", 3},
		{"With punctuation", "Hello, World!", 4},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.EstimateTokens(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d tokens, got %d for input: %q", tt.expected, result, tt.input)
			}
		})
	}
}

func TestGPTEstimate(t *testing.T) {
	tokenizer := New(ModelGPT4)
	
	tests := []struct {
		name     string
		input    string
		minTokens int
		maxTokens int
	}{
		{"Empty string", "", 0, 0},
		{"Single word", "hello", 1, 1},
		{"Two words", "hello world", 2, 2},
		{"With punctuation", "Hello, world!", 4, 5},
		{"Long word", "internationalization", 5, 6},
		{"Numbers", "12345", 2, 3},
		{"Mixed content", "The year 2024 is here!", 7, 9},
		{"Unicode", "Hello 世界", 3, 4},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.EstimateTokens(tt.input)
			if result < tt.minTokens || result > tt.maxTokens {
				t.Errorf("Expected between %d and %d tokens, got %d for input: %q", 
					tt.minTokens, tt.maxTokens, result, tt.input)
			}
		})
	}
}

func TestClaudeEstimate(t *testing.T) {
	tokenizer := New(ModelClaude)
	
	tests := []struct {
		name     string
		input    string
		minTokens int
		maxTokens int
	}{
		{"Empty string", "", 0, 0},
		{"Single word", "hello", 1, 1},
		{"Two words", "hello world", 2, 2},
		{"With punctuation", "Hello, world!", 4, 5},
		{"Long sentence", "The quick brown fox jumps over the lazy dog.", 9, 12},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.EstimateTokens(tt.input)
			if result < tt.minTokens || result > tt.maxTokens {
				t.Errorf("Expected between %d and %d tokens, got %d for input: %q", 
					tt.minTokens, tt.maxTokens, result, tt.input)
			}
		})
	}
}

func TestTokenizeVerbose(t *testing.T) {
	tokenizer := New(ModelGPT4)
	
	input := "Hello, world! 123"
	result := tokenizer.TokenizeVerbose(input)
	
	if len(result) == 0 {
		t.Error("Expected verbose output, got empty result")
	}
	
	totalTokens := 0
	for _, info := range result {
		totalTokens += info.Tokens
		if info.Text == "" {
			t.Error("Token info should not have empty text")
		}
		if info.Type == "" {
			t.Error("Token info should have a type")
		}
	}
	
	if totalTokens == 0 {
		t.Error("Total tokens should be greater than 0")
	}
}

func BenchmarkSimpleEstimate(b *testing.B) {
	tokenizer := New(ModelSimple)
	text := "The quick brown fox jumps over the lazy dog. This is a sample text for benchmarking."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.EstimateTokens(text)
	}
}

func BenchmarkGPTEstimate(b *testing.B) {
	tokenizer := New(ModelGPT4)
	text := "The quick brown fox jumps over the lazy dog. This is a sample text for benchmarking."
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokenizer.EstimateTokens(text)
	}
}