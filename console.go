package gox

import (
	"bytes"
	"fmt"
	"log"
)

const (
	ConsoleColorPrefixFirst  = "\x1b["
	ConsoleColorPrefixSecond = ";1m"
	ConsoleColorSuffix       = "\x1b[0m"
	ConsoleColorYellow       = "33"
	ConsoleColorBrightBlue   = "94"
	ConsoleColorBlue         = "34"
	ConsoleColorRed          = "31"
	ConsoleColorGreen        = "32"
)

func LogConsoleWriteLabeledValueI(label string, value interface{})  {
	log.Println(ConsoleWriteLabeledValueI(label, value))
}

func ConsoleWriteLabeledValueI(label string, value interface{}) string {
	bb := &bytes.Buffer{}

	bb.WriteString(ConsoleColorPrefixFirst)
	bb.WriteString("34")
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(label)
	bb.WriteString(": ")
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString(ConsoleColorPrefixFirst)
	bb.WriteString("91")
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(fmt.Sprintf("%v", value))
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString("\n")

	return bb.String()
}

func ConsoleWriteLabeledValue(label string, value string) string {
	bb := &bytes.Buffer{}

	bb.WriteString(ConsoleColorPrefixFirst)
	bb.WriteString("34")
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(label)
	bb.WriteString(": ")
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString(ConsoleColorPrefixFirst)
	bb.WriteString("91")
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(value)
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString("\n")

	return bb.String()
}

func ConsoleWriteLabeledValueWithCheck(label string, value string, verified bool) string {
	bb := &bytes.Buffer{}

	bb.WriteString(ConsoleColorPrefixFirst)
	bb.WriteString("34")
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(label)
	bb.WriteString(": ")
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString(ConsoleColorPrefixFirst)
	if verified {
		bb.WriteString("32")
	} else {
		bb.WriteString("31")
	}
	bb.WriteString(ConsoleColorPrefixSecond)
	bb.WriteString(value)
	bb.WriteString(ConsoleColorSuffix)
	bb.WriteString("\n")

	return bb.String()
}

func ConsoleTestColors() {
	for i := 0; i < 100; i++ {

		bb := &bytes.Buffer{}

		bb.WriteString(ConsoleColorPrefixFirst)
		bb.WriteString(fmt.Sprintf("%v", i))
		bb.WriteString(ConsoleColorPrefixSecond)
		bb.WriteString(fmt.Sprintf("Color: %v", i))
		bb.WriteString(ConsoleColorSuffix)
		fmt.Println(bb)
	}
}
