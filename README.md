# tokencount

A fast command-line tool to estimate token counts for AI models (GPT-3.5, GPT-4, Claude, etc.).

## Features

- 🚀 Fast token estimation without external dependencies
- 📊 Support for multiple AI model token patterns
- 🌍 Cross-platform (Linux, macOS, Windows)
- 📝 Multiple input methods (direct text, files, stdin)
- 🔍 Verbose mode with detailed token breakdown
- 📋 JSON output for programmatic use
- 🎯 Lightweight and minimal binary size

## Installation

### Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/wearewebera/tokencount/releases).

#### macOS/Linux
```bash
# Download and extract (replace with your platform and version)
tar -xzf tokencount_Darwin_x86_64.tar.gz

# Make executable
chmod +x tokencount

# Move to PATH (optional)
sudo mv tokencount /usr/local/bin/
```

#### Windows
Download the Windows zip file, extract it, and run `tokencount.exe` from the command prompt.

### Build from Source

```bash
git clone https://github.com/wearewebera/tokencount.git
cd tokencount
go build -o tokencount
```

## Usage

### Basic Usage

```bash
# Estimate tokens for direct text
tokencount "Hello, world!"

# Estimate tokens from a file
tokencount -f document.txt

# Pipe text from another command
echo "Some text" | tokencount
cat document.txt | tokencount
```

### Model Selection

The tool supports different estimation algorithms optimized for various AI models:

```bash
# Use GPT-4 estimation (default)
tokencount "Your text here"

# Use GPT-3.5 estimation
tokencount -m gpt-3.5 "Your text here"

# Use Claude estimation
tokencount -m claude "Your text here"

# Use simple estimation (4 chars = 1 token)
tokencount -m simple "Your text here"
```

### Output Options

```bash
# Verbose output with token breakdown
tokencount -v "Hello, world!"

# JSON output for programmatic use
tokencount -j "Your text here"

# Combine options
tokencount -m gpt-4 -v -j "Detailed analysis"
```

### Examples

#### Simple Text
```bash
$ tokencount "The quick brown fox jumps over the lazy dog"
Tokens: 9
Model: gpt-4 (GPT-4 estimation algorithm)
```

#### Verbose Mode
```bash
$ tokencount -v "Hello, world! 123"
Tokens: 5
Model: gpt-4 (GPT-4 estimation algorithm)

Token breakdown:
--------------------------------------------------
Hello                word         1 token(s)
,                    punctuation  1 token(s)
world                word         1 token(s)
!                    punctuation  1 token(s)
123                  number       1 token(s)

Original text:
--------------------------------------------------
Hello, world! 123
```

#### JSON Output
```bash
$ echo "API request" | tokencount -j
{
  "token_count": 2,
  "model": "gpt-4",
  "model_info": "GPT-4 estimation algorithm"
}
```

#### File Input
```bash
$ tokencount -f large_document.txt
Tokens: 1523
Model: gpt-4 (GPT-4 estimation algorithm)
```

## How It Works

The tool uses different algorithms to estimate token counts based on the selected model:

- **Simple**: Basic estimation using 4 characters = 1 token rule
- **GPT-3.5/GPT-4**: Advanced algorithm considering:
  - Word boundaries and length
  - Punctuation (typically 1 token each)
  - Numbers (grouped by digits)
  - Unicode characters (1 token per character for CJK, etc.)
- **Claude**: Similar to GPT but with Claude-specific optimizations

Note: These are estimations. Actual token counts may vary slightly from the real tokenizers used by AI providers.

## Command Line Options

```
Usage:
  tokencount [options] [text]
  echo "text" | tokencount [options]
  tokencount -f file.txt [options]

Options:
  -m, --model string     Model to use for estimation (simple, gpt-3.5, gpt-4, claude) (default "gpt-4")
  -f, --file string      Input file to read
  -v, --verbose          Verbose output with token breakdown
  -j, --json             Output results as JSON
  -h, --help             Show help
  --version              Show version information
```

## Development

### Running Tests
```bash
go test ./... -v
```

### Building for Multiple Platforms
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o tokencount-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o tokencount-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o tokencount-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o tokencount-windows-amd64.exe
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.