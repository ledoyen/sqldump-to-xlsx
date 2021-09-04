package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	excelize "github.com/xuri/excelize/v2"
)

func main() {
	os.Exit(checkArgAndTransformFile())
}

func checkArgAndTransformFile() int {
	if len(os.Args) < 2 {
		log.Fatal("An input file must be provided")
		return 1
	}
	filePath := os.Args[1]
	return transformFile(filePath)
}

func transformFile(filePath string) int {
	start := time.Now()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer file.Close()

	workbook := writesLinesToWorkbook(file)

	err2 := workbook.SaveAs(filePath + ".xlsx")
	if err2 != nil {
		fmt.Println(err2)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Converted in", elapsed)
	return 0
}

type workbookState struct {
	tableName  string
	sheetIndex int
	rowIndex   int
}

func writesLinesToWorkbook(file *os.File) *excelize.File {
	workbook := excelize.NewFile()

	boldStyle, err := workbook.NewStyle(`{"font": {"bold": true}}`)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)

	state := workbookState{}

	for {
		line, err := r.ReadString('\n')
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				break
			}
		}
		line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
		writeLineToWorkbook(workbook, &state, line, boldStyle)
	}

	workbook.DeleteSheet("Sheet1")

	return workbook
}

var /* const */ linePattern = regexp.MustCompile(`^INSERT INTO (?P<Table>[^\s]+) \((?P<Cols>.+)\) VALUES \((?P<Values>.+)\);$`)

func writeLineToWorkbook(workbook *excelize.File, state *workbookState, line string, boldStyle int) {
	tokens := linePattern.FindStringSubmatch(line)
	if tokens != nil {
		writeHeadersIfFirstInsert(workbook, state, tokens[1], tokens[2], boldStyle)
		state.rowIndex++
		writeLineValues(workbook, state, tokens[3])
	}
}

func writeHeadersIfFirstInsert(workbook *excelize.File, state *workbookState, tableName string, rawHeaders string, headerStyle int) {
	if tableName == state.tableName {
		return
	}
	state.tableName = tableName
	state.sheetIndex = workbook.NewSheet(state.tableName)
	state.rowIndex = 1
	headers := strings.Split(rawHeaders, ",")
	for i, header := range headers {
		coords, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			panic(err)
		}
		error := workbook.SetCellStr(state.tableName, coords, header)
		if error != nil {
			fmt.Printf("Unable to set value for cell %s: %s\n", coords, error)
		}
	}
	colName, err := excelize.ColumnNumberToName(len(headers))
	if err != nil {
		panic(err)
	}
	error := workbook.SetCellStyle(state.tableName, "A1", "A"+colName, headerStyle)
	if error != nil {
		fmt.Printf("Unable to set style for headers: %s\n", error)
	}
}

func writeLineValues(workbook *excelize.File, state *workbookState, rawValues string) {
	for i, value := range split(rawValues) {
		coords, err := excelize.CoordinatesToCellName(i+1, state.rowIndex)
		if err != nil {
			panic(err)
		}
		error := workbook.SetCellStr(state.tableName, coords, limit(value))
		if error != nil {
			fmt.Printf("Unable to set value for cell %s: %s\n", coords, error)
		}
	}
}

func split(v string) []string {
	var inQuote = false
	var tokens []string
	var sb strings.Builder
	for pos, c := range v {
		if c == '\\' && pos < len(v) && v[pos+1] == '\'' {
			sb.WriteRune('\'')
		} else if c == '\'' && pos > 0 && v[pos-1] == '\\' {
			continue
		} else if c == '\'' {
			inQuote = !inQuote
		} else if c == ',' && !inQuote {
			tokens = append(tokens, strings.TrimSpace(sb.String()))
			sb = strings.Builder{}
		} else {
			sb.WriteRune(c)
		}
	}
	if sb.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(sb.String()))
	}

	return tokens
}

func limit(s string) string {
	if len(s) > 32763 {
		a := []rune(s)
		return string(a[0:32763])
	}
	return s
}
