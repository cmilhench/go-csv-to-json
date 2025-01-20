package json

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cm-dev/template/internal/core/time"
)

func Flatten(data any) map[string]any {
	out := make(map[string]any)
	flatten(data, "", out)
	return out
}

func flatten(data any, prefix string, result map[string]any) {
	switch value := data.(type) {
	case map[string]any:
		for k, v := range value {
			flatten(v, fmt.Sprintf("%s.%s", prefix, k), result)
		}
	case []any:
		for i, v := range value {
			flatten(v, fmt.Sprintf("%s.%d", prefix, i), result)
		}
	default:
		result[strings.TrimPrefix(prefix, ".")] = value
	}
}

func Unflatten(data map[string]any) (any, error) {
	var (
		object = make(map[string]any, 0)
		array  = make([]any, 0)
		nested = make(map[string]map[string]any, 0)
	)

	for key, value := range data {
		parts := strings.Split(key, ".")
		if len(parts) == 1 {
			if _, err := strconv.Atoi(parts[0]); err == nil {
				array = append(array, value)
				continue
			}
			object[parts[0]] = value
			continue
		}
		if _, exists := nested[parts[0]]; !exists {
			nested[parts[0]] = make(map[string]any, 1)
		}
		nested[parts[0]][strings.Join(parts[1:], ".")] = value
	}

	for key, value := range nested {
		if _, err := strconv.Atoi(key); err == nil {
			if inner, err := Unflatten(value); err == nil {
				array = append(array, inner)
			}
			continue
		}
		inner, err := Unflatten(value)
		if err != nil {
			return nil, fmt.Errorf("couldn't unflatten inner: %s", err)
		}
		object[key] = inner
	}

	if len(array) > 0 {
		return array, nil
	}
	return object, nil
}

func Inflate(headers, row []string) (any, error) {
	data := make(map[string]any)
	for i, header := range headers {
		if len(row) > i {
			data[header] = parseValue(row[i])
		}
	}
	unflattened, err := Unflatten(data)
	if err != nil {
		return nil, fmt.Errorf("error inflating data: %v", err)
	}
	return unflattened, nil
}

func parseValue(value string) any {
	// Try converting to an integer
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	// Try converting to a float
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		return floatValue
	}
	// Convert boolean values
	if value == "true" || value == "false" {
		return value == "true"
	}
	// Try converting to a time
	if floatValue, err := time.Parse(value, time.TimeFormats...); err == nil {
		return floatValue
	}
	// Return as string if nothing else matches
	return value
}
