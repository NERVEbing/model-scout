package cli

import (
	"fmt"
	"strings"

	"github.com/NERVEbing/model-scout/internal/platform"
)

type filterExpressions []string

func (f *filterExpressions) String() string {
	return strings.Join(*f, ",")
}

func (f *filterExpressions) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type filterOp int

const (
	filterEq filterOp = iota
	filterNotEq
)

type filter struct {
	key    string
	op     filterOp
	values []string
}

func parseFilters(inputs []string) ([]filter, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	parsed := make([]filter, 0, len(inputs))
	for _, raw := range inputs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		var (
			key   string
			value string
			op    filterOp
		)
		if strings.Contains(raw, "!=") {
			parts := strings.SplitN(raw, "!=", 2)
			key = strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[1])
			op = filterNotEq
		} else if strings.Contains(raw, "=") {
			parts := strings.SplitN(raw, "=", 2)
			key = strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[1])
			op = filterEq
		} else {
			return nil, fmt.Errorf("invalid filter %q (expected key=value or key!=value)", raw)
		}

		key = strings.ToLower(key)
		switch key {
		case "available", "status", "model", "platform":
		default:
			return nil, fmt.Errorf("unsupported filter key: %s", key)
		}

		values := splitFilterValues(value)
		if len(values) == 0 {
			return nil, fmt.Errorf("invalid filter %q (missing value)", raw)
		}
		if key == "available" {
			for _, entry := range values {
				if _, err := parseBool(entry); err != nil {
					return nil, err
				}
			}
		}
		parsed = append(parsed, filter{key: key, op: op, values: values})
	}
	return parsed, nil
}

func splitFilterValues(raw string) []string {
	parts := strings.Split(raw, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		values = append(values, trimmed)
	}
	return values
}

func applyFilters(results []platform.ProbeResult, filters []filter) []platform.ProbeResult {
	filtered := make([]platform.ProbeResult, 0, len(results))
	for _, result := range results {
		if matchesAllFilters(result, filters) {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func matchesAllFilters(result platform.ProbeResult, filters []filter) bool {
	for _, f := range filters {
		if !matchesFilter(result, f) {
			return false
		}
	}
	return true
}

func matchesFilter(result platform.ProbeResult, f filter) bool {
	switch f.key {
	case "available":
		value := result.Available
		return matchBool(value, f)
	case "status":
		return matchString(result.Status, f)
	case "model":
		return matchString(result.Model, f)
	case "platform":
		return matchString(result.Platform, f)
	default:
		return false
	}
}

func matchBool(actual bool, f filter) bool {
	for _, raw := range f.values {
		value, err := parseBool(raw)
		if err != nil {
			return false
		}
		if actual == value {
			return f.op == filterEq
		}
	}
	return f.op == filterNotEq
}

func matchString(actual string, f filter) bool {
	for _, value := range f.values {
		if actual == value {
			return f.op == filterEq
		}
	}
	return f.op == filterNotEq
}

func parseBool(raw string) (bool, error) {
	switch strings.ToLower(raw) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean %q", raw)
	}
}
