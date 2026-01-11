package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NERVEbing/model-scout/internal/output"
	"github.com/NERVEbing/model-scout/internal/platform"
	"github.com/NERVEbing/model-scout/internal/platform/dashscope"
	"github.com/NERVEbing/model-scout/internal/scout"
)

const defaultDashscopeKeyEnv = "DASHSCOPE_API_KEY"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "scan":
		if err := runScan(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "usage: model-scout scan [flags]")
}

func runScan(args []string) error {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	platformName := flags.String("platform", "", "platform to scan")
	apiKey := flags.String("api-key", "", "api key")
	keyEnv := flags.String("key-env", defaultDashscopeKeyEnv, "environment variable for api key")
	workers := flags.Int("workers", 4, "number of workers")
	timeout := flags.Duration("timeout", 15*time.Second, "http timeout")
	outFormat := flags.String("out", "json", "output format: json or yaml")
	outputFile := flags.String("output-file", "", "output file path")
	onlyOK := flags.Bool("only-ok", false, "output only available models")
	exclude := flags.String("exclude", "", "comma-separated substrings to exclude")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if *platformName == "" {
		return errors.New("--platform is required")
	}

	key := strings.TrimSpace(*apiKey)
	if key == "" {
		key = strings.TrimSpace(os.Getenv(*keyEnv))
	}
	if key == "" {
		return fmt.Errorf("api key missing; provide --api-key or set %s", *keyEnv)
	}

	platformImpl, err := platformFromName(*platformName, key, *timeout)
	if err != nil {
		return err
	}

	engine := scout.Engine{Platform: platformImpl, Workers: *workers}
	ctx := context.Background()
	excludes := splitExclude(*exclude)
	results, err := engine.Scan(ctx, excludes)
	if err != nil {
		return err
	}

	if *onlyOK {
		filtered := results[:0]
		for _, result := range results {
			if result.Status == "ok" {
				filtered = append(filtered, result)
			}
		}
		results = filtered
	}

	return writeOutput(*outFormat, *outputFile, results)
}

func splitExclude(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		filtered = append(filtered, trimmed)
	}
	return filtered
}

func platformFromName(name, apiKey string, timeout time.Duration) (platform.Platform, error) {
	switch strings.ToLower(name) {
	case "dashscope":
		return dashscope.NewPlatform(apiKey, timeout), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", name)
	}
}

func writeOutput(format, outputFile string, payload []platform.ProbeResult) error {
	format = strings.ToLower(format)
	var err error
	var writer *os.File
	if outputFile == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(outputFile)
		if err != nil {
			return err
		}
		defer writer.Close()
	}

	switch format {
	case "json":
		return output.WriteJSON(writer, payload)
	case "yaml":
		return output.WriteYAML(writer, payload)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
