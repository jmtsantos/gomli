package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	// notations
	RgxClassName    = regexp.MustCompile(`\.class\s(.*)`)
	RgxFields       = regexp.MustCompile(`\.field\s+(.*)`)
	RgxMethodsStart = regexp.MustCompile(`\.method\s+(.*)`)
	RgxMethodsEnd   = regexp.MustCompile(`\.end method`)
	RgxMethodName   = regexp.MustCompile(`\s+(.*)\(`)

	// instructions
	RgxInstruction0x1A = regexp.MustCompile(`(const-string) (.*), "(.*)"`)
	RgxInstruction0x26 = regexp.MustCompile(`(fill-array-data) (.*), (.*)`)
	RgxInstruction0x71 = regexp.MustCompile(`(invoke-static) {((?:.*,?))}, (.*)->(.*)`)
)

type Smali struct {
	Path      string
	ClassName string
	Raw       string
	RawSlc    []string

	Properties   []string
	Methods      []Method
	Instructions []Instruction
}

type Method struct {
	Name               string
	LineStart, LineEnd int
}

type Instruction struct {
	Raw        string
	LineNumber int

	Verbs []string

	OpCode byte
	Method string
	Class  string
}

type Call struct {
	Class  string
	Method string
}

func (s *Smali) ParseClassName() {
	resultSlc := RgxClassName.FindStringSubmatch(s.Raw)

	if len(strings.Split(resultSlc[1], " ")) > 1 {
		classNameSplit := strings.Split(resultSlc[1], " ")
		s.ClassName = classNameSplit[len(classNameSplit)-1]
	} else {
		s.ClassName = resultSlc[1]
	}
}

func (s *Smali) ParseProperties() {
	resultSlc := RgxFields.FindAllString(s.Raw, -1)

	s.Properties = append(s.Properties, resultSlc...)
}

// Reads instructions inside a method
func (s *Smali) ParseMethods() {

	resultSlc := strings.Split(s.Raw, "\n")

	addInstruction := false
	methodName := ""
	for lineNumber, line := range resultSlc {
		line = strings.TrimSpace(line)

		if line == "" {
			s.Instructions = append(s.Instructions, Instruction{
				Raw:        "",
				LineNumber: lineNumber,
			})
		}

		if RgxMethodsStart.MatchString(line) {
			addInstruction = true
			methodName = RgxMethodsStart.FindStringSubmatch(line)[1]

			if len(strings.Split(methodName, " ")) > 1 {
				methodNameSplit := strings.Split(methodName, " ")
				methodName = methodNameSplit[len(methodNameSplit)-1]

			} else {
				methodName = resultSlc[1]
			}

			s.Methods = append(s.Methods, Method{
				Name:      methodName,
				LineStart: lineNumber,
			})

		}
		if RgxMethodsEnd.MatchString(line) {
			addInstruction = false
			s.Methods[len(s.Methods)-1].LineEnd = lineNumber

		}

		if addInstruction {
			s.Instructions = append(s.Instructions, s.ParseInstruction(methodName, lineNumber, line))
		} else {
			s.Instructions = append(s.Instructions, s.ParseInstruction("", lineNumber, line))
		}

	}
}

func (s *Smali) ParseInstruction(methodName string, lineNbr int, line string) (inst Instruction) {
	// http://pallergabor.uw.hu/androidblog/dalvik_opcodes.html
	inst.Raw = line
	inst.Verbs = strings.Split(line, " ")
	inst.LineNumber = lineNbr
	inst.Method = methodName
	inst.Class = s.ClassName

	switch inst.Verbs[0] {

	// 00 nop
	case "nop":
		inst.OpCode = 0x00

	// 1A const-string vx,string_id
	case "const-string":
		inst.OpCode = 0x1A
		inst.Verbs = RgxInstruction0x1A.FindStringSubmatch(line)

		// 26 fill-array-data v4, :array_a8
	case "fill-array-data":
		inst.OpCode = 0x26
		inst.Verbs = RgxInstruction0x26.FindStringSubmatch(line)

	// 71 invoke-static {parameters}, methodtocall
	case "invoke-static":
		inst.OpCode = 0x71
		inst.Verbs = RgxInstruction0x71.FindStringSubmatch(line)
	}

	return
}

func (s *Smali) GenerateRAWSmali() {
	var joinedInstructions string
	for _, v := range s.Instructions {
		joinedInstructions = fmt.Sprintf("%s%s\n", joinedInstructions, v.Raw)
	}
	s.Raw = joinedInstructions
}

// SaveSmali overwrites the smali code with the patched version
func (s *Smali) SaveSmali() (err error) {
	f, err := os.OpenFile(s.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	f.WriteString(s.Raw)
	if err := f.Close(); err != nil {
		return err
	}
	f.Close()
	return
}
