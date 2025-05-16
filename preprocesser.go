package main

import (
	"bufio"
	"fmt"
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
	if len(lineParts) < 3 {
		return []string{"invalid mv instruction"}
	}

	rd := strings.TrimSpace(lineParts[1])
	rs := strings.TrimSpace(lineParts[2])

	return []string{fmt.Sprintf("addi %s, %s, 0", rd, rs)}
}

func handleJ(lineParts []string) []string {
	if len(lineParts) < 2 {
		return []string{"invalid j instruction"}
	}

	target := strings.TrimSpace(lineParts[1])
	return []string{fmt.Sprintf("jal x0, %s", target)}
}

func handleJR(lineParts []string) []string {
	if len(lineParts) < 2 {
		return []string{"invalid jr instruction"}
	}

	rs := strings.TrimSpace(lineParts[1])
	return []string{fmt.Sprintf("jalr x0, %s, 0", rs)}
}

func handleBLE(lineParts []string) []string {
	if len(lineParts) < 4 {
		return []string{"invalid ble instruction"}
	}

	rs1 := strings.TrimSpace(lineParts[1])
	rs2 := strings.TrimSpace(lineParts[2])
	label := strings.TrimSpace(lineParts[3])

	// ble is equivalent to bge with swapped operands
	return []string{fmt.Sprintf("bge %s, %s, %s", rs2, rs1, label)}
}

func handleLI(lineParts []string) []string {
	if len(lineParts) < 3 {
		return []string{"invalid li instruction"}
	}

	rd := strings.TrimSpace(lineParts[1])
	imm := strings.TrimSpace(lineParts[2])

	val, err := strconv.Atoi(imm)
	if err != nil {
		//imm may be a register
		return []string{
			fmt.Sprintf("lui %s, %%hi(%s)", rd, imm),
			fmt.Sprintf("addi %s, %s, %%lo(%s)", rd, rd, imm),
		}
	}

	if val >= -2048 && val < 2048 {
		return []string{fmt.Sprintf("addi %s, x0, %s", rd, imm)}
	}

	return []string{
		fmt.Sprintf("lui %s, %%hi(%s)", rd, imm),
		fmt.Sprintf("addi %s, %s, %%lo(%s)", rd, rd, imm),
	}
}

func handleLA(lineParts []string) []string {
	if len(lineParts) < 3 {
		return []string{"invalid la instruction"}
	}

	rd := strings.TrimSpace(lineParts[1])
	symbol := strings.TrimSpace(lineParts[2])

	return []string{
		fmt.Sprintf("auipc %s, %%pcrel_hi(%s)", rd, symbol),
		fmt.Sprintf("addi %s, %s, %%pcrel_lo(%s)", rd, rd, symbol),
	}
}
