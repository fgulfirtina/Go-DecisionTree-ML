package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

// reading the csv file into dataset [][]string
func readCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	dataset, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return dataset
}

// calculating entropy for a list of class labels
func entropy(labels []string) float64 {
	count := make(map[string]int)
	for _, label := range labels {
		count[label]++
	}
	total := float64(len(labels))

	sum := 0.0
	for _, c := range count {
		p := float64(c) / total
		sum -= p * math.Log2(p)
	}
	return sum
}

// calculating information gain for a specific attribute index
func information_gain(rows [][]string, attributeIndex int) float64 {
	subsets := make(map[string][][]string)
	for _, row := range rows {
		key := row[attributeIndex]
		subsets[key] = append(subsets[key], row)
	}
	total := float64(len(rows))
	weightedEntropy := 0.0

	for _, subset := range subsets {
		labels := []string{}
		for _, row := range subset {
			labels = append(labels, row[len(row)-1])
		}
		w := float64(len(subset)) / total
		weightedEntropy += w * entropy(labels)
	}

	mainLabels := []string{}
	for _, row := range rows {
		mainLabels = append(mainLabels, row[len(row)-1])
	}
	return entropy(mainLabels) - weightedEntropy
}

// finds the best attribute to split on
func find_best_attribute(rows [][]string, header []string) int {
	bestGain := -1.0
	bestIndex := -1

	for i := 0; i < len(header)-1; i++ {
		gain := information_gain(rows, i)
		fmt.Printf("Information gain for %s: %.4f\n", header[i], gain)
		if gain > bestGain {
			bestGain = gain
			bestIndex = i
		}
	}
	return bestIndex
}

// returns the majority class label in a list of rows
func majority_label(rows [][]string) string {
	count := make(map[string]int)
	for _, row := range rows {
		label := row[len(row)-1]
		count[label]++
	}
	major := ""
	max := -1
	for k, v := range count {
		if v > max {
			max = v
			major = k
		}
	}
	return major
}

// recursively builds the decision tree
func build_tree(rows [][]string, header []string) interface{} {
	labels := []string{}
	for _, row := range rows {
		labels = append(labels, row[len(row)-1])
	}

	first := labels[0]
	allSame := true
	for _, label := range labels {
		if label != first {
			allSame = false
			break
		}
	}
	if allSame {
		return first
	}

	if len(rows[0]) == 1 {
		return majority_label(rows)
	}

	bestIndex := find_best_attribute(rows, header)
	bestAttr := header[bestIndex]
	tree := map[string]interface{}{bestAttr: map[string]interface{}{}}

	seen := map[string]bool{}
	for _, row := range rows {
		seen[row[bestIndex]] = true
	}

	for val := range seen {
		subset := [][]string{}
		for _, row := range rows {
			if row[bestIndex] == val {
				newRow := append([]string{}, row[:bestIndex]...)
				newRow = append(newRow, row[bestIndex+1:]...)
				subset = append(subset, newRow)
			}
		}
		subHeader := append([]string{}, header[:bestIndex]...)
		subHeader = append(subHeader, header[bestIndex+1:]...)
		subtree := build_tree(subset, subHeader)
		tree[bestAttr].(map[string]interface{})[val] = subtree
	}
	return tree
}

// prints the decision tree
func print_tree(tree interface{}, indent string) {
	switch t := tree.(type) {
	case string:
		fmt.Println(indent + "Prediction: " + t)
	case map[string]interface{}:
		for attr, branches := range t {
			for val, subtree := range branches.(map[string]interface{}) {
				fmt.Printf("%s%s: %s ->\n", indent, attr, val)
				print_tree(subtree, indent+"  ")
			}
		}
	}
}

// walks the tree using user input
func predict(tree interface{}, header []string, input []string) string {
	for {
		t, ok := tree.(map[string]interface{})
		if !ok {
			return tree.(string)
		}
		attr := ""
		for k := range t {
			attr = k
			break
		}
		index := -1
		for i, h := range header {
			if h == attr {
				index = i
				break
			}
		}
		val := input[index]
		next, exists := t[attr].(map[string]interface{})[val]
		if !exists {
			return "Unknown"
		}
		tree = next
	}
}

// asks the user for input and shows prediction
func prediction(tree interface{}, header []string) {
	fmt.Println("\nEnter attribute values for prediction (press Enter to quit):")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		input := []string{}
		for i := 0; i < len(header)-1; i++ {
			fmt.Printf("%s: ", header[i])
			scanner.Scan()
			val := strings.TrimSpace(scanner.Text())
			if val == "" {
				return
			}
			input = append(input, val)
		}
		result := predict(tree, header, input)
		fmt.Printf("Prediction: %s\n\n", result)
	}
}

func main() {
	dataset := readCSV("dataset.csv")
	header := dataset[0]
	data := dataset[1:]

	tree := build_tree(data, header)
	fmt.Println()
	print_tree(tree, "")
	prediction(tree, header)
}
