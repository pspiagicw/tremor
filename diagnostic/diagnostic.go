package diagnostic

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pspiagicw/tremor/token"
)

type Span struct {
	StartOffset int
	EndOffset   int
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
}

type Diagnostic struct {
	Stage   string
	Message string
	File    string
	Source  string
	Span    *Span
}

func New(stage string, file string, source string, format string, args ...any) *Diagnostic {
	return &Diagnostic{
		Stage:   stage,
		Message: fmt.Sprintf(format, args...),
		File:    file,
		Source:  source,
	}
}

func NewAtToken(stage string, file string, source string, tok *token.Token, width int, format string, args ...any) *Diagnostic {
	d := New(stage, file, source, format, args...)
	d.Span = SpanFromToken(tok, width)
	return d
}

func SpanFromToken(tok *token.Token, width int) *Span {
	if tok == nil {
		return nil
	}
	if width < 1 {
		width = 1
	}

	startLine := tok.Line
	startColumn := tok.Column
	if startLine < 1 {
		startLine = 1
	}
	if startColumn < 1 {
		startColumn = 1
	}

	return &Span{
		StartOffset: tok.Offset,
		EndOffset:   tok.Offset + width,
		StartLine:   startLine,
		EndLine:     startLine,
		StartColumn: startColumn,
		EndColumn:   startColumn + width - 1,
	}
}

func (d *Diagnostic) Error() string {
	return d.Message
}

func (d *Diagnostic) Pretty() string {
	if d == nil {
		return ""
	}

	header := "error"
	if d.Stage != "" {
		header = fmt.Sprintf("error[%s]", d.Stage)
	}

	useColor := colorEnabled()
	headerLabel := header
	if useColor {
		headerLabel = style(header, ansiBold, ansiRed)
	}
	if d.Message == "" {
		return headerLabel
	}

	if d.Span == nil || d.Source == "" {
		return fmt.Sprintf("%s: %s", headerLabel, d.Message)
	}

	file := d.File
	if file == "" {
		file = "<input>"
	}

	lineNo := d.Span.StartLine
	col := d.Span.StartColumn
	if lineNo < 1 {
		lineNo = 1
	}
	if col < 1 {
		col = 1
	}

	sourceLine := getLine(d.Source, lineNo)
	gutterWidth := len(strconv.Itoa(lineNo))
	caretWidth := d.Span.EndColumn - d.Span.StartColumn + 1
	if caretWidth < 1 {
		caretWidth = 1
	}
	if col > len(sourceLine)+1 {
		col = len(sourceLine) + 1
	}

	var b strings.Builder
	arrow := "-->"
	separator := "|"
	underline := strings.Repeat("^", caretWidth)
	if useColor {
		arrow = style(arrow, ansiBlue)
		separator = style(separator, ansiBlue)
		underline = style(underline, ansiRed)
	}

	fmt.Fprintf(&b, "%s: %s\n", headerLabel, d.Message)
	fmt.Fprintf(&b, " %s %s:%d:%d\n", arrow, file, lineNo, col)
	fmt.Fprintf(&b, "%*s %s\n", gutterWidth, "", separator)
	fmt.Fprintf(&b, "%*d %s %s\n", gutterWidth, lineNo, separator, sourceLine)
	fmt.Fprintf(&b, "%*s %s %s%s", gutterWidth, "", separator, strings.Repeat(" ", col-1), underline)
	return b.String()
}

func Render(err error) string {
	if err == nil {
		return ""
	}
	if d, ok := err.(*Diagnostic); ok {
		return d.Pretty()
	}
	return err.Error()
}

func getLine(source string, line int) string {
	if line < 1 {
		return ""
	}

	currentLine := 1
	start := 0
	for i := 0; i < len(source); i++ {
		if source[i] == '\n' {
			if currentLine == line {
				return source[start:i]
			}
			currentLine += 1
			start = i + 1
		}
	}
	if currentLine == line {
		return source[start:]
	}
	return ""
}

const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"
	ansiRed   = "\x1b[31m"
	ansiBlue  = "\x1b[34m"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func colorEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	switch strings.ToLower(os.Getenv("TREMOR_COLOR")) {
	case "0", "false", "off", "never":
		return false
	case "1", "true", "on", "always":
		return true
	}

	term := strings.ToLower(os.Getenv("TERM"))
	if term == "" || term == "dumb" {
		return false
	}

	return true
}

func style(input string, codes ...string) string {
	if input == "" || len(codes) == 0 {
		return input
	}
	return strings.Join(codes, "") + input + ansiReset
}

func StripANSI(input string) string {
	return ansiPattern.ReplaceAllString(input, "")
}
