package main

import (
	"fmt"
	"io"
	"strings"
)

type TableWriter struct {
	w      io.Writer
	rows   [][]string
	header []string
	maxes  []int
}

const (
	Center = "+"
	Row    = "-"
	Column = "|"
)

func NewTableWriter(w io.Writer) *TableWriter {
	return &TableWriter{
		w: w,
	}
}

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
	w.printLine()
	w.print(w.header)
	w.printLine()
	for _, row := range w.rows {
		w.print(row)
	}
	w.printLine()
}

func (w *TableWriter) printLine() {
	fmt.Fprintf(w.w, "+")
	for _, max := range w.maxes {
		fmt.Fprintf(w.w, "%s", strings.Repeat("-", max))
		fmt.Fprintf(w.w, "%s", "+")
	}
	fmt.Fprintf(w.w, "\n")
}

func (w *TableWriter) print(row []string) {
	fmt.Fprintf(w.w, "|")
	for i, val := range row {
		spacesCnt := w.maxes[i] - len(val) - 1
		fmt.Fprintf(w.w, " %s%s", val, strings.Repeat(" ", spacesCnt))
		fmt.Fprintf(w.w, "%s", "|")
	}
	fmt.Fprintf(w.w, "\n")
}
