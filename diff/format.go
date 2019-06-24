package diff

// TODO: Calculate whether an added map/slice is all new (>>>) or has changes (|||)

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aws-cloudformation/rain/format"
)

const indent = "  "

func Format(d Diff) string {
	switch v := d.(type) {
	case diffSlice:
		return formatSlice(v)
	case diffMap:
		return formatMap(v)
	case diffValue:
		f := format.NewFormatter()
		f.SetCompact()
		return f.Format(v.value)
	case mode:
		if v == Unchanged {
			return "No changes\n"
		}

		return fmt.Sprintf("%sEverything!\n", v)
	}

	panic(fmt.Sprintf("Unexpected %#v\n", d))
}

func stubValue(v diffValue) string {
	switch v.value.(type) {
	case map[string]interface{}:
		return "{...}"
	case []interface{}:
		return "[...]"
	default:
		return "..."
	}
}

func formatSlice(d diffSlice) string {
	output := strings.Builder{}

	for i, v := range d {
		m := v.mode()

		if m == Unchanged {
			continue
		}

		output.WriteString(fmt.Sprintf("%s[%d]:", m.String(), i))

		if m == Removed {
			output.WriteString(" " + stubValue(v.(diffValue)) + "\n")
		} else {
			output.WriteString(formatSub(v))
		}
	}

	return output.String()
}

func formatMap(d diffMap) string {
	output := strings.Builder{}

	// Sort the keys
	keys := make([]string, 0)

	for k, _ := range d {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := d[k]
		m := v.mode()

		if m == Unchanged {
			continue
		}

		output.WriteString(fmt.Sprintf("%s%s:", m.String(), k))

		if m == Removed {
			output.WriteString(" " + stubValue(v.(diffValue)) + "\n")
		} else {
			output.WriteString(formatSub(v))
		}
	}

	return output.String()
}

func formatSub(d Diff) string {
	// Format the element
	formatted := strings.Split(Format(d), "\n")

	v, isValue := d.(diffValue)
	if isValue {
		k := reflect.ValueOf(v.value).Kind()

		if k != reflect.Array && k != reflect.Map && k != reflect.Slice {
			// It's a scalar
			return fmt.Sprintf(" %s\n", formatted[0])
		}
	} else if len(formatted) == 1 {
		// It's a scalar
		return fmt.Sprintf(" %s\n", formatted[0])
	}

	// Trim out blank lines
	parts := make([]string, 0)
	for _, part := range formatted {
		if strings.TrimSpace(part) != "" {
			parts = append(parts, part)
		}
	}

	output := strings.Builder{}

	output.WriteString("\n")
	for _, part := range parts {
		if isValue {
			part = Added.String() + indent + part
		} else {
			part = part[:len(Added.String())] + indent + part[len(Added.String()):]
		}

		output.WriteString(part)
		output.WriteString("\n")
	}

	return output.String()
}
