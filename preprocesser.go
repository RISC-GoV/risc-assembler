package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Preprocess(file *os.File) []string {

	var result []string = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineIndx := strings.Index(line, "#")
		if lineIndx != -1 {
			line = line[:lineIndx]
		}
		line = strings.ReplaceAll(line, "\t", " ")

		lineParts := strings.Split(line, " ")
		lineParts = removeEmptyStrings(lineParts)

		if len(lineParts) == 0 {
			continue
		}

		if res, ok := PseudoToInstruction[lineParts[0]]; ok {
			var resArray []string = res(lineParts)

			for _, res := range resArray {
				result = append(result, res)
			}
		} else {
			result = append(result, line)
		}

	}

	return result
}

func handleMV(lineParts []string) []string {
	// split by, Then add instruction as last element
	if len(lineParts) > 2 {
		lineParts = []string{lineParts[0], strings.Join(lineParts[1:], "")}
	}
	lineParts = append(strings.Split(lineParts[1], ","), lineParts[0])
	return []string{"addi " + lineParts[0] + "," + lineParts[1] + ",0"}
}

func handleJ(lineParts []string) []string {
	return []string{"jal x0," + lineParts[1]}
}

func handleJR(lineParts []string) []string {
	return []string{"jalr x0," + lineParts[1] + ",0"}
}

func handleBLE(lineParts []string) []string {
	// split by, Then add instruction as last element
	if len(lineParts) > 2 {
		lineParts = []string{lineParts[0], strings.Join(lineParts[1:], "")}
	}
	lineParts = append(strings.Split(lineParts[1], ","), lineParts[0])
	return []string{"bge " + lineParts[1] + "," + lineParts[0] + "," + lineParts[2]}
}

func handleLI(lineParts []string) []string {
	var result []string = make([]string, 0)
	if len(lineParts) > 2 {
		lineParts = []string{lineParts[0], strings.Join(lineParts[1:], "")}
	}
	// split by, Then add instruction as last element
	lineParts = append(strings.Split(lineParts[1], ","), lineParts[0])
	res, err := strconv.Atoi(lineParts[len(lineParts)-2])
	if err != nil {
		panic(err)
	}
	if res >= -2048 && res < 2048 {
		result = append(result, "addi "+lineParts[0]+",x0,"+lineParts[1])
	} else {
		result = append(result, "lui "+lineParts[0]+",%hi("+lineParts[1]+")")
		result = append(result, "addi "+lineParts[0]+","+lineParts[0]+",%lo("+lineParts[1]+")")
	}
	return result
}
