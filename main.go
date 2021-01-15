package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	args := os.Args
	cmd := filepath.Base(args[0])
	if len(args) < 2 {
		fmt.Printf("Usage: %s files.csv|files.xlsx\n", cmd)
		return
	}
	for _, filename := range args[1:] {
		ext := strings.ToLower(filepath.Ext(filename))
		src, err := os.Open(filename)
		check(err)
		switch ext {
		case ".csv":
			csvToXlsx(src, filename)
		case ".xlsx":
			xlsxToCsv(src, filename)
		}
	}
}

func csvToXlsx(src io.Reader, filename string) {
	csvFile := csv.NewReader(src)
	xlsxFile := excelize.NewFile()
	row := 0
	for {
		record, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		check(err)
		row++
		for k, v := range record {
			xlsxFile.SetCellValue("Sheet1", axis(k, row), v)
		}
	}
	err := xlsxFile.SaveAs(filename[:len(filename)-4] + ".xlsx")
	check(err)
}

func xlsxToCsv(src io.Reader, filename string) {
	xlsxFile, err := excelize.OpenReader(src)
	check(err)
	for _, sheet := range xlsxFile.GetSheetMap() {
		name := filename[:len(filename)-5]
		if sheet != "Sheet1" {
			name += "_" + sheet
		}
		csvFile, err := os.Create(name + ".csv")
		check(err)
		csvWriter := csv.NewWriter(csvFile)
		for _, row := range xlsxFile.GetRows(sheet) {
			err = csvWriter.Write(row)
			check(err)
		}
		csvWriter.Flush()
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func axis(x, y int) string {
	return fmt.Sprintf("%s%d", col(x), y)
}

func col(n int) string {
	if n < 26 {
		return string(n + 'A')
	}
	return col(n/26-1) + string(n%26+'A')
}
