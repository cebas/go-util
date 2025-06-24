package util

import "fmt"

type Log struct {
	verbosity int
	indent    int
}

func NewLog() Log {
	return Log{verbosity: 0, indent: 0}
}

func (l *Log) Verbosity(verbosity int) {
	l.verbosity = verbosity
}

func _print(msg string) {
	fmt.Print(msg)
}

func _println(msg string) {
	fmt.Println(msg)
}

func _printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

func (l *Log) printIndent() {
	for i := 0; i < l.indent; i++ {
		_print("\t")
	}
}

func (l *Log) Indent(newLevel int) {
	l.indent = newLevel
}

func (l *Log) shouldPrint(verbose int) bool {
	return verbose <= l.verbosity
}

// Println prints a message to stdout
func (l *Log) Println(verbose int, msg string) {
	if l.shouldPrint(verbose) {
		l.printIndent()
		_println(msg)
	}
}

// Printf prints a formatted message to stdout
func (l *Log) Printf(verbose int, format string, v ...any) {
	if l.shouldPrint(verbose) {
		l.printIndent()
		_printf(format, v...)
	}
}

// Logln prints a message to stdout
// Deprecated: use Log.Println instead
func Logln(msg string) {
	_println(msg)
}

// Logf prints a formatted message to stdout
// Deprecated: use Log.Printf instead
func Logf(format string, v ...any) {
	_printf(format, v...)
}
