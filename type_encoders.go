package assembler

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func TranslateRType(opcode int, rd int, func3 int, rs1 int, rs2 int, func7 int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	func7 &= 0b1111111

	//shifts
	res := uint32(func7)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return res
}

func TranslateIType(opcode int, rd int, func3 int, rs1 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	imm &= 0b111111111111

	//shifts
	res := uint32(imm)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return res
}

func TranslateSType(opcode int, func3 int, rs1 int, rs2 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	imm4_0 := imm & 0b11111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm11_5 := imm
	imm11_5 >>= 5
	imm11_5 &= 0b1111111

	//shifts
	res := uint32(imm11_5)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 5
	res |= uint32(imm4_0)
	res <<= 7
	res |= uint32(opcode)

	return res
}

func TranslateBType(opcode int, func3 int, rs1 int, rs2 int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	imm11 := imm >> 11 & 0b1
	imm4_1 := imm >> 1 & 0b1111
	func3 &= 0b111
	rs1 &= 0b11111
	rs2 &= 0b11111
	imm10_5 := imm >> 5 & 0b111111
	imm12 := imm >> 12 & 0b1

	//shifts
	res := uint32(imm12)
	res <<= 6
	res |= uint32(imm10_5)
	res <<= 5
	res |= uint32(rs2)
	res <<= 5
	res |= uint32(rs1)
	res <<= 3
	res |= uint32(func3)
	res <<= 4
	res |= uint32(imm4_1)
	res <<= 1
	res |= uint32(imm11)
	res <<= 7
	res |= uint32(opcode)

	return res
}

func TranslateUType(opcode int, rd int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm &= 0b11111111111111111111

	//shifts
	res := uint32(imm)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return res
}
func TranslateJType(opcode int, rd int, imm int) uint32 {
	//masks
	opcode &= 0b1111111
	rd &= 0b11111
	imm19_12 := imm >> 12 & 0b11111111
	imm11 := imm >> 11 & 0b1
	imm10_1 := imm >> 1 & 0b1111111111
	var imm20 int
	if imm >= 0 {
		imm20 = 0
	} else {
		imm20 = 1
	}

	//shifts
	res := uint32(imm20)
	res <<= 10
	res |= uint32(imm10_1)
	res <<= 1
	res |= uint32(imm11)
	res <<= 8
	res |= uint32(imm19_12)
	res <<= 5
	res |= uint32(rd)
	res <<= 7
	res |= uint32(opcode)

	return res
}

func (p *Program) InstructionToBinary(t *Token, relativeInstrCount int) (uint32, error) {
	if t.tokenType != instruction {
		return 0, errors.New("expected instruction")
	}
	if t.value == "ecall" || t.value == "ebreak" {
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rd := 0
		rs1 := 0
		func12 := int(t.opPair.opByte[2])

		return TranslateIType(opcode, rd, func3, rs1, func12), nil
	}
	switch t.opPair.opType {
	case R: // add x0, x0, x0
		if len(t.children) != 3 {
			return 0, errors.New(t.value + " is not a valid instruction")
		}
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		func7 := int(t.opPair.opByte[2])
		rd, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		rs1, err := t.children[1].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		rs2, err := t.children[2].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		return TranslateRType(opcode, rd, func3, rs1, rs2, func7), nil
	case I: // lw x0, 0(x0)
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rd, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		rs1, imm, err := p.parseComplexValue(t.children[1], relativeInstrCount)
		if err != nil {
			return 0, err
		}
		if len(t.children) > 2 {
			//if two children and rs1 = 0 then we can guess that previous one was a register (addi X5, X5, 5)
			if rs1 == 0 {
				rs1 = imm
			}

			_, imm, err = p.parseComplexValue(t.children[2], relativeInstrCount)
			if err != nil {
				return 0, err
			}
		}
		imm |= (int(t.opPair.opByte[2]) << 5)
		return TranslateIType(opcode, rd, func3, rs1, imm), nil
	case S: //sw x0, 0(x0)
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rs1, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		rs2, imm, err := p.parseComplexValue(t.children[1], relativeInstrCount)
		if err != nil {
			return 0, err
		}
		return TranslateSType(opcode, func3, rs1, rs2, imm), nil
	case B: // beq x0, x0, 0
		opcode := int(t.opPair.opByte[0])
		func3 := int(t.opPair.opByte[1])
		rs1, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		rs2, err := t.children[1].getRegisterNumericValue()
		tmp, imm, err := p.parseComplexValue(t.children[2], relativeInstrCount)
		if err != nil {
			return 0, err
		}
		if imm == 0 {
			imm = tmp
		}
		return TranslateBType(opcode, func3, rs1, rs2, imm), nil
	case U: // lui x0, 0
		opcode := int(t.opPair.opByte[0])
		rd, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}

		tmp, imm, err := p.parseComplexValue(t.children[1], relativeInstrCount)
		if err != nil {
			return 0, err
		}
		if imm == 0 {
			imm = tmp
		}
		return TranslateUType(opcode, rd, imm), nil
	case J: // jal x0, 0
		opcode := int(t.opPair.opByte[0])
		rd, err := t.children[0].getRegisterNumericValue()
		if err != nil {
			return 0, err
		}
		tmp, imm, err := p.parseComplexValue(t.children[1], relativeInstrCount)
		if err != nil {
			return 0, err
		}
		if imm == 0 {
			imm = tmp
		}
		return TranslateJType(opcode, rd, imm), nil

		//COMPRESSED INSTRUCTIONS HAVE BEEN IMPLEMENTED VIA PREPROCESSER AS THE CPU IS NOT
		//SUPPOSED TO BE CAPABLE OF HANDLING COMPRESSED INSTRUCTIONS
	case CI:
		return 0, errors.New("todo: compressed instructions (CI)") //todo
	case CSS:
		return 0, errors.New("todo: compressed instructions (CSS)") //todo
	case CL:
		return 0, errors.New("todo: compressed instructions (CL)") //todo
	case CJ:
		return 0, errors.New("todo: compressed instructions (CJ)") //todo
	case CR:
		return 0, errors.New("todo: compressed instructions (CR)") //todo
	case CB:
		return 0, errors.New("todo: compressed instructions (CB)") //todo
	case CIW:
		return 0, errors.New("todo: compressed instructions (CIW)") //todo
	case CS:
		return 0, errors.New("todo: compressed instructions (CS)") //todo
	default:
		return 0, errors.New("unhandled default case")
	}
}

// Updated parseIntValue function to handle hex values and character literals
func parseIntValue(valueStr string) (int, error) {
	// Strip whitespace
	valueStr = strings.TrimSpace(valueStr)

	// Handle character literal 'X'
	if len(valueStr) >= 3 && valueStr[0] == '\'' && valueStr[len(valueStr)-1] == '\'' {
		// Extract the character between the quotes
		char := valueStr[1 : len(valueStr)-1]

		// Handle escape sequences if present
		if len(char) > 1 && char[0] == '\\' {
			switch char[1] {
			case 'n':
				return int('\n'), nil
			case 'r':
				return int('\r'), nil
			case 't':
				return int('\t'), nil
			case '\\':
				return int('\\'), nil
			case '\'':
				return int('\''), nil
			case '0':
				return 0, nil
				//we could add more but too long
			default:
				return 0, fmt.Errorf("unsupported escape sequence: %s", char)
			}
		} else if len(char) == 1 {
			// Single character
			return int(char[0]), nil
		} else {
			return 0, fmt.Errorf("invalid character literal: %s", valueStr)
		}
	}

	base := 10
	// Check if this is a hex value
	if strings.HasPrefix(valueStr, "-0x") || strings.HasPrefix(valueStr, "-0X") {
		valueStr = "-" + valueStr[3:] // Remove the 0x prefix
		base = 16
	} else if strings.HasPrefix(valueStr, "0x") || strings.HasPrefix(valueStr, "0X") {
		valueStr = valueStr[2:] // Remove the 0x prefix
		base = 16
	}

	// First parse as int64 to handle the full range of 32-bit unsigned values
	val, err := strconv.ParseInt(valueStr, base, 64)
	if err != nil {
		return 0, err
	}

	// Check if value is within 32-bit range
	if val < math.MinInt32 || val > math.MaxUint32 {
		return 0, fmt.Errorf("value %s is out of range for 32-bit architecture", valueStr)
	}
	return int(int32(val)), nil
}

func (p *Program) parseComplexValue(tok *Token, relativeInstrCount int) (int, int, error) {
	switch tok.tokenType { //0(x0) OU 0(.LC1)
	case complexValue:
		if tok.children[0].tokenType == literal {
			if tok.children[1].tokenType == register {
				reg, err := p.parseLabelOrLiteral(tok.children[1], relativeInstrCount)
				if err != nil {
					return 0, 0, err
				}
				imm, err := parseIntValue(tok.children[0].value)
				if err != nil {
					return 0, 0, err
				}
				return reg, int(imm), nil
			} else { //constant
				con, err := p.parseLabelOrLiteral(tok.children[1], relativeInstrCount)
				if err != nil {
					return 0, 0, errors.New(tok.children[1].value + " not found")
				}
				imm, err := parseIntValue(tok.children[0].value)
				if err != nil {
					return 0, 0, err
				}
				return con, int(imm), nil
			}
		} else {
			parsed, err := p.parseLabelOrLiteral(tok.children[1], relativeInstrCount)
			if err != nil {
				return 0, 0, err
			}
			parsed, err = handleModifier(tok.children[0].value, parsed, relativeInstrCount)
			if err != nil {
				return 0, 0, err
			}
			return 0, parsed, nil
		}
	case varLabel:
		fallthrough
	case varValue:
		fallthrough
	case constantValue:
		fallthrough
	case literal:
		fallthrough
	case register:
		val, err := p.parseLabelOrLiteral(tok, relativeInstrCount)
		if err != nil {
			return 0, 0, err
		}
		return 0, val, nil
	}
	return 0, 0, nil
}

func handleModifier(mod string, val int, relativeInstrCount int) (int, error) {
	switch mod {
	case "%lo":
		return (val + relativeInstrCount) & 0xFFF, nil //int(uint8(val)), nil
	case "%hi":
		return ((val + relativeInstrCount) >> 12) & 0xFFFFF, nil
	case "%pcrel_lo":
		return val & 0xFFF, nil //int(uint8(val)), nil
	case "%pcrel_hi":
		return (val >> 12) & 0xFFFFF, nil
	}
	return 0, errors.New("modifier not found")
}

func (p *Program) parseLabelOrLiteral(tok *Token, instructionRelativePos int) (int, error) {
	switch tok.tokenType {
	case varLabel:
		fallthrough
	case varValue:
		fallthrough
	case constantValue:
		imm, ok := p.compilationVariables.labelPositions[tok.value]
		if !ok {
			return 0, errors.New(tok.value + " not found")
		}
		imm -= instructionRelativePos
		return imm, nil
	case literal:
		val, err := parseIntValue(tok.value)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	case register:
		imm, err := matchTokenValid(tok.value)
		if err != nil {
			return 0, err
		}
		return imm, nil
	}
	return 0, errors.New(fmt.Sprintf("wrong token type:  %+v", tok))
}
