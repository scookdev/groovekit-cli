package output

import (
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Colors
var (
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()
	Error   = color.New(color.FgRed, color.Bold).SprintFunc()
)

// Table is a wrapper around go-pretty table
type Table struct {
	writer table.Writer
}

// NewTable creates a nicely formatted table
func NewTable(headers []string) *Table {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)

	// Convert headers to table.Row
	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	return &Table{
		writer: t,
	}
}

// Render is a no-op for compatibility (go-pretty renders on Flush)
func (t *Table) Render() {
	// No-op: go-pretty handles rendering automatically
}

// Append adds a row to the table
func (t *Table) Append(row []string) {
	// Convert string slice to table.Row
	tableRow := make(table.Row, len(row))
	for i, cell := range row {
		tableRow[i] = cell
	}
	t.writer.AppendRow(tableRow)
}

// Flush writes all buffered data to output
func (t *Table) Flush() {
	t.writer.Render()
}

// SuccessMessage prints a green success message
func SuccessMessage(msg string) {
	color.Green("✓ " + msg)
}

// ErrorMessage prints a red error message
func ErrorMessage(msg string) {
	color.Red("✗ " + msg)
}

// InfoMessage prints a cyan info message
func InfoMessage(msg string) {
	color.Cyan(msg)
}
