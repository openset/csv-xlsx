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

const VERSION = "0.1.2"

func main() {
	args := os.Args
	cmd := filepath.Base(args[0])
	if len(args) < 2 {
		fmt.Printf("Usage: %s file.csv|file.xlsx...\n", cmd)
		return
	}
	for _, filename := range args[1:] {
		src, err := os.Open(filename)
		check(err)
		ext := strings.ToLower(filepath.Ext(filename))
		name := strings.TrimSuffix(filename, ext)
		basename := filepath.Base(filename)
		switch ext {
		case ".csv":
			csvToXlsx(src, name)
		case ".xlsx":
			xlsxToCsv(src, name)
		default:
			fmt.Printf("skipping file %s\n", basename)
		}
		_ = src.Close()
	}
}

func csvToXlsx(src io.Reader, name string) {
	maxRow := 0xFFFFF + 1
	csvFile := csv.NewReader(src)
	csvFile.FieldsPerRecord = -1
	xlsxFile := excelize.NewFile()
	for row := 0; ; row++ {
		record, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		check(err)
		sheet := fmt.Sprintf("Sheet%d", row/maxRow+1)
		if row > 0 && row%maxRow == 0 {
			xlsxFile.NewSheet(sheet)
		}
		for k, v := range record {
			xlsxFile.SetCellValue(sheet, axis(k, row%maxRow+1), v)
		}
	}
	xlsxFile.SetActiveSheet(1)
	err := xlsxFile.SaveAs(name + ".xlsx")
	check(err)
}

func xlsxToCsv(src io.Reader, name string) {
	xlsxFile, err := excelize.OpenReader(src)
	check(err)
	for _, sheet := range xlsxFile.GetSheetMap() {
		filename := name
		if sheet != "Sheet1" {
			filename += "_" + sheet
		}
		csvFile, err := os.Create(filename + ".csv")
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

func axis(col, row int) string {
	return fmt.Sprintf("%s%d", convertToTitle(col), row)
}

func convertToTitle(n int) string {
	if n < 26 {
		return fmt.Sprintf("%c", n+'A')
	}
	return convertToTitle(n/26-1) + fmt.Sprintf("%c", n%26+'A')
}
