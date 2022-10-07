package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type TableWriter struct {
	Writer io.Writer
	rows   [][]string
	header []string
	maxes  []int
}

const (
	Center = "+"
	Row    = "-"
	Column = "|"
)

func (w *TableWriter) SetHeader(hdr []string) {
	w.maxes = make([]int, len(hdr))
	for i := range hdr {
		w.maxes[i] = len(hdr[i]) + 2
	}
	w.header = hdr
}

func (w *TableWriter) Append(row []string) {
	for i := range row {
		if w.maxes[i] < len(row[i])+2 {
			w.maxes[i] = len(row[i]) + 2
		}
	}
	w.rows = append(w.rows, row)
}

func (w *TableWriter) Render() {
	if w.Writer == nil {
		w.Writer = os.Stdout // default
	}
	w.printLine()
	w.print(w.header)
	w.printLine()
	for _, row := range w.rows {
		w.print(row)
	}
	w.printLine()
}

func (w *TableWriter) printLine() {
	fmt.Fprintf(w.Writer, "+")
	for _, max := range w.maxes {
		fmt.Fprintf(w.Writer, "%s", strings.Repeat("-", max))
		fmt.Fprintf(w.Writer, "%s", "+")
	}
	fmt.Fprintf(w.Writer, "\n")
}

func (w *TableWriter) print(row []string) {
	fmt.Fprintf(w.Writer, "|")
	for i, val := range row {
		spacesCnt := w.maxes[i] - len(val) - 1
		fmt.Fprintf(w.Writer, " %s%s", val, strings.Repeat(" ", spacesCnt))
		fmt.Fprintf(w.Writer, "%s", "|")
	}
	fmt.Fprintf(w.Writer, "\n")
}
