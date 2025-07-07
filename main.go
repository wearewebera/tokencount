package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/wearewebera/tokencount/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var opts cmd.Options
	var showVersion bool
	var showHelp bool

	flag.StringVar(&opts.Model, "m", "gpt-4", "Model to use for estimation (simple, gpt-3.5, gpt-4, claude)")
	flag.StringVar(&opts.Model, "model", "gpt-4", "Model to use for estimation (simple, gpt-3.5, gpt-4, claude)")
	flag.StringVar(&opts.File, "f", "", "Input file to read")
	flag.StringVar(&opts.File, "file", "", "Input file to read")
	flag.BoolVar(&opts.Verbose, "v", false, "Verbose output with token breakdown")
	flag.BoolVar(&opts.Verbose, "verbose", false, "Verbose output with token breakdown")
	flag.BoolVar(&opts.JsonOutput, "j", false, "Output results as JSON")
	flag.BoolVar(&opts.JsonOutput, "json", false, "Output results as JSON")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.BoolVar(&showHelp, "help", false, "Show help")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "tokencount - Estimate token count for AI models\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  tokencount [options] [text]\n")
		fmt.Fprintf(os.Stderr, "  echo \"text\" | tokencount [options]\n")
		fmt.Fprintf(os.Stderr, "  tokencount -f file.txt [options]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  tokencount \"Hello, world!\"\n")
		fmt.Fprintf(os.Stderr, "  tokencount -m gpt-3.5 -v \"Detailed analysis\"\n")
		fmt.Fprintf(os.Stderr, "  echo \"Some text\" | tokencount -j\n")
		fmt.Fprintf(os.Stderr, "  tokencount -f document.txt -m claude\n\n")
		fmt.Fprintf(os.Stderr, "Available models:\n")
		fmt.Fprintf(os.Stderr, "  simple   - Basic estimation (4 chars = 1 token)\n")
		fmt.Fprintf(os.Stderr, "  gpt-3.5  - GPT-3.5 estimation\n")
		fmt.Fprintf(os.Stderr, "  gpt-4    - GPT-4 estimation (default)\n")
		fmt.Fprintf(os.Stderr, "  claude   - Claude estimation\n")
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("tokencount version %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Get text from command line arguments
	text := strings.Join(flag.Args(), " ")

	// Execute the command
	if err := cmd.Execute(text, opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}