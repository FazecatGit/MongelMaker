package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ExportRecord struct {
	Timestamp string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int64
	RSI       *float64
	ATR       *float64
	Analysis  string
	Signals   []string
}

func ExportRecordsToCSV(filename string, bars []ExportRecord) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"Timestamp", "Open", "High", "Low", "Close", "Volume", "RSI", "ATR", "Analysis", "Signals"}); err != nil {
		return err
	}

	for _, bar := range bars {

		if err := writer.Write(RecordToRow(bar)); err != nil {
			return err
		}
	}

	return nil
}

func RecordToRow(record ExportRecord) []string {
	row := []string{
		record.Timestamp,
		strconv.FormatFloat(record.Open, 'f', 2, 64),
		strconv.FormatFloat(record.High, 'f', 2, 64),
		strconv.FormatFloat(record.Low, 'f', 2, 64),
		strconv.FormatFloat(record.Close, 'f', 2, 64),
		strconv.FormatInt(record.Volume, 10),
	}
	if record.RSI != nil {
		row = append(row, strconv.FormatFloat(*record.RSI, 'f', 2, 64))
	} else {
		row = append(row, "")
	}
	if record.ATR != nil {
		row = append(row, strconv.FormatFloat(*record.ATR, 'f', 2, 64))
	} else {
		row = append(row, "")
	}
	row = append(row, record.Analysis)
	row = append(row, strings.Join(record.Signals, "; "))
	return row
}

func ExportRecordsToJSON(filename string, records []ExportRecord) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(records)
}

func ExportData(format, filename string, records []ExportRecord) error {
	filename = "exported_data/" + filename
	switch format {
	case "csv":
		return ExportRecordsToCSV(filename, records)
	case "json":
		return ExportRecordsToJSON(filename, records)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
