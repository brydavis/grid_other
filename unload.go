package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/tealeg/xlsx"
)

func (p Panel) Write(fname string) { // ADD ERROR
	Unload(p, fname)
}

func (p Panel) Dump(fname string) { // ADD ERROR
	Unload(p, fname)
}

func (p Panel) Unload(fname string) { // ADD ERROR
	Unload(p, fname)
}

func Unload(p Panel, fname string) { // ADD ERROR
	// p.Clean()
	switch path.Ext(fname) {
	case ".csv":
		// Convert to CSV-like structure
		width := p.Size().Width
		length := p.Size().Length
		headers := make(map[string]int)
		colnum := 0
		for head, _ := range p {
			headers[head] = colnum
			colnum += 1
		}

		records := [][]string{}
		record := make([]string, width)
		for head, _ := range p {
			record[headers[head]] = head
		}
		records = append(records, record)

		for i := 0; i < length; i++ {
			record := make([]string, width)
			for head, col := range p {
				var val string
				switch t := col[i].(type) {
				case time.Time:
					d110, _ := dateFormat(110)
					tt := t.Format(d110)
					fmt.Println(tt)
					val = tt
				default:
					val = fmt.Sprintf("%v", t)
				}
				record[headers[head]] = val
			}
			records = append(records, record)
		}

		// fmt.Println(records)
		// Write to file
		file := new(os.File)
		if _, err := os.Stat(fname); err == nil {
			file, _ = os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0777)
		} else {
			file, _ = os.Create(fname)
		}

		w := csv.NewWriter(file)
		// w.WriteAll(records)

		for _, record := range records {
			if err := w.Write(record); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		w.Flush() // Write any buffered data to the underlying writer
		if err := w.Error(); err != nil {
			log.Fatal(err)
		}

	case ".xml":

	case ".tsv":

	case ".json":
		export := []map[string]interface{}{}

		// width := p.Size().Width
		length := p.Size().Length

		for i := 0; i < length; i++ {
			row := make(map[string]interface{})
			for header, series := range p {
				row[header] = series[i]
			}
			export = append(export, row)
		}

		j, _ := json.Marshal(export)

		file := new(os.File)
		if _, err := os.Stat(fname); err == nil {
			file, _ = os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0777)
		} else {
			file, _ = os.Create(fname)
		}

		_, err := file.Write(j)
		if err != nil {
			fmt.Println(err)
		}

		file.Close()

	case ".xlsx", ".xls":
		// Convert to CSV-like structure
		width := p.Size().Width
		length := p.Size().Length
		headers := make(map[string]int)
		colnum := 0
		for head, _ := range p {
			headers[head] = colnum
			colnum += 1
		}

		records := [][]string{}
		record := make([]string, width)
		for head, _ := range p {
			record[headers[head]] = head
		}
		records = append(records, record)

		for i := 0; i < length; i++ {
			record := make([]string, width)
			for head, col := range p {
				var val string
				switch t := col[i].(type) {
				case time.Time:
					d110, _ := dateFormat(110)
					tt := t.Format(d110)
					fmt.Println(tt)
					val = tt
				default:
					val = fmt.Sprintf("%v", t)
				}
				record[headers[head]] = val
			}
			records = append(records, record)
		}

		// fmt.Println(records)
		// Write to file
		// file := new(os.File)
		// if _, err := os.Stat(fname); err == nil {
		// 	file, _ = os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0777)
		// } else {
		// 	file, _ = os.Create(fname)
		// }

		var excel *xlsx.File
		var sheet *xlsx.Sheet
		var row *xlsx.Row
		var cell *xlsx.Cell
		var err error

		excel = xlsx.NewFile()
		sheet, err = excel.AddSheet("Sheet1")
		if err != nil {
			fmt.Printf(err.Error())
		}

		for _, rec := range records {
			row = sheet.AddRow()

			for _, field := range rec {
				cell = row.AddCell()
				cell.Value = fmt.Sprintf("%s", field)
			}
		}

		err = excel.Save(fname)
		if err != nil {
			fmt.Printf(err.Error())
			log.Fatalln("error writing record to xlsx:", err)

		}

	default:
		file := new(os.File)
		if _, err := os.Stat(fname); err == nil {
			file, _ = os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0777)
		} else {
			file, _ = os.Create(fname)
		}

		file.Write([]byte(fmt.Sprintf("%s", p)))
	}
}
