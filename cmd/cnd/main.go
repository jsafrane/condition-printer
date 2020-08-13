package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v3"
	operatorv1 "github.com/openshift/api/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

type Operator struct {
	Status operatorv1.OperatorStatus `json:"status"`
}

type operatorList struct {
	Items []Operator `json:"items"`
}

// Split input by "}" or "---"
func scanYAMLJSON(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n}\n")); i >= 0 {
		i = i + 3
		return i, data[0:i], nil
	}
	if i := bytes.Index(data, []byte("\n---\n")); i >= 0 {
		i = i + 5
		return i, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func isJSON(data []byte) bool {
	data = bytes.TrimLeft(data, "\n ")
	return len(data) == 0 || data[0] == '{'
}

func parseJSON(data []byte) (Operator, error) {
	var op Operator
	if bytes.Contains(data, []byte(`"kind": "List"`)) {
		list := operatorList{}
		err := json.Unmarshal(data, &list)
		if err != nil {
			return op, err
		}
		if len(list.Items) == 0 {
			return op, fmt.Errorf("No object on input")
		}
		return list.Items[0], nil
	}

	err := json.Unmarshal(data, &op)
	return op, err
}

func parseYAML(data []byte) (Operator, error) {
	var op Operator
	if bytes.Contains(data, []byte(`kind: List`)) {
		list := operatorList{}
		err := yaml.Unmarshal(data, &list)
		if err != nil {
			return op, err
		}
		if len(list.Items) == 0 {
			return op, fmt.Errorf("No object on input")
		}
		return list.Items[0], nil
	}

	err := yaml.Unmarshal(data, &op)
	return op, err
}

func parseConditions(data []byte) ([]operatorv1.OperatorCondition, error) {
	var op Operator
	var err error

	if isJSON(data) {
		op, err = parseJSON(data)
	} else {
		op, err = parseYAML(data)
	}
	if err != nil {
		return nil, err
	}

	return op.Status.Conditions, nil
}

func clearScreen() {
	os.Stdout.Write([]byte("\033[H\033[2J"))
}

func printConditions(conditions []operatorv1.OperatorCondition) {
	clearScreen()
	time := findLastTimestamp(conditions)
	fmt.Printf("Last transition time: %s\n\n", time)

	var t = make(table, 0)
	for _, cnd := range conditions {
		var r row
		r = append(r, item{aurora.WhiteFg, cnd.Type})
		r = append(r, item{conditionColor(cnd), string(cnd.Status)})
		r = append(r, item{aurora.WhiteFg, cnd.Message})
		t = append(t, r)
	}
	t.Sort()
	t.Print()
}

func findLastTimestamp(conditions []operatorv1.OperatorCondition) string {
	var last metav1.Time
	for _, cnd := range conditions {
		if last.IsZero() {
			last = cnd.LastTransitionTime
		}
		if last.Before(&cnd.LastTransitionTime) {
			last = cnd.LastTransitionTime
		}
	}
	return last.String()
}

func conditionColor(cnd operatorv1.OperatorCondition) aurora.Color {
	type colorSet struct {
		trueColor, falseColor aurora.Color
	}
	colors := map[string]colorSet{
		"Available": {
			aurora.GreenFg + aurora.BoldFm,
			aurora.YellowFg + aurora.BoldFm,
		},
		"Progressing": {
			aurora.YellowFg + aurora.ItalicFm,
			aurora.GreenFg + aurora.ItalicFm,
		},
		"Degraded": {
			aurora.RedFg + aurora.UnderlineFm,
			aurora.GreenFg + aurora.UnderlineFm,
		},
	}

	for suffix, set := range colors {
		if strings.HasSuffix(cnd.Type, suffix) {
			switch cnd.Status {
			case operatorv1.ConditionFalse:
				return set.falseColor
			case operatorv1.ConditionTrue:
				return set.trueColor
			}
		}
	}
	// The default color is white
	return aurora.WhiteFg
}

func main() {
	var scanner *bufio.Scanner

	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scanner = bufio.NewScanner(f)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	scanner.Split(scanYAMLJSON)
	for scanner.Scan() {
		cnds, err := parseConditions(scanner.Bytes())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		printConditions(cnds)
	}
	if scanner.Err() != nil {
		// handle error.
		panic(scanner.Err())
	}
}
