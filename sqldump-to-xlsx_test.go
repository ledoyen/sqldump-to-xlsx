package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	excelize "github.com/xuri/excelize/v2"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestSplit(t *testing.T) {
	type Case struct {
		rawValue       string
		expectedValues []string
	}

	for _, tt := range []Case{
		{"a value", []string{"a value"}},
		{"a, b", []string{"a", "b"}},
		{"'a','b,c',d", []string{"a", "b,c", "d"}},
	} {
		t.Run(fmt.Sprintf("%s %+v", "Splitting", tt.rawValue), func(t *testing.T) {
			values := split(tt.rawValue)

			require.Equal(t, tt.expectedValues, values)
		})
	}
}

func TestLimit(t *testing.T) {
	type Case struct {
		inputSize    int
		expectedSize int
	}
	for _, tt := range []Case{
		{7, 7},
		{25485, 25485},
		{45870, 32763},
	} {
		t.Run(fmt.Sprintf("%s %+v", "Limiting text of size", tt.inputSize), func(t *testing.T) {

			text := strings.Repeat("A", tt.inputSize)
			outputText := limit(text)

			require.Len(t, outputText, tt.expectedSize)
		})
	}
}

func TestE2E(t *testing.T) {
	inputFile := "testdata/dump1.sql"
	outputFile := inputFile + ".xlsx"
	exitCode := transformFile(inputFile)
	require.Equal(t, 0, exitCode)
	workbook, err := excelize.OpenFile(outputFile)
	require.NoError(t, err, "Failed to read the expected XLSX ouput file %s", outputFile)

	type Case struct {
		sheet         string
		cell          string
		expectedValue string
	}

	for _, tt := range []Case{
		{"user", "A1", "id"},
		{"user", "C1", "age"},
		{"user", "B2", "George"},
		{"user", "A3", "2"},
		{"user", "C3", "32"},
		{"room", "B4", "Living room"},
	} {
		t.Run(fmt.Sprintf("%s %s", tt.sheet, tt.cell), func(t *testing.T) {
			cellValue, err := workbook.GetCellValue(tt.sheet, tt.cell)
			require.NoError(t, err, "Failed to read Cell %s of sheet %s", tt.cell, tt.sheet)
			require.Equal(t, tt.expectedValue, cellValue)
		})
	}
}
