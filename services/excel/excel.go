package excel

import (
	"bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"strconv"
)

var Columns = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ"}

type Value interface{}

type Headers []string

type Excel struct {
	excel *excelize.File
	index int
}

type Row []Value

type Rows []Row

func New() *Excel {
	return &Excel{
		index: 2,
		excel: excelize.NewFile(),
	}
}

func Axis(col, row int) string {
	return TransferCellName(col, row)
}

// TransferCellName returns like A1, A2, BB1, col started with 1.
func TransferCellName(col, row int) string {
	var cols string
	v := col
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		cols = string(rune(k+64)) + cols
	}
	if row > 0 {
		return cols + strconv.Itoa(row)
	}
	return cols
}

func WriteRow(f *excelize.File, sheet string, startCol int, row int, data []interface{}, styleId ...int) {
	for _, val := range data {
		f.SetCellValue(sheet, TransferCellName(startCol, row), val)
		if len(styleId) > 0 {
			f.SetCellStyle(sheet, TransferCellName(startCol, row), TransferCellName(startCol, row), styleId[0])
		}
		startCol += 1
	}
}

func (receiver *Excel) SetHeader(headers Headers) {
	for k, v := range headers {
		column := Columns[k]
		receiver.excel.SetCellValue("Sheet1", column+strconv.Itoa(1), v)
	}
}

func (receiver *Excel) AddRow(row Row) {
	for j, column := range row {
		columnIndex := Columns[j]
		receiver.excel.SetCellValue("Sheet1", columnIndex+strconv.Itoa(receiver.index), column)
	}
	receiver.index++
}

func (receiver *Excel) SetRows(rows Rows) {
	for k, row := range rows {
		rowIndex := strconv.Itoa(k + 1 + 1)
		for j, column := range row {
			columnIndex := Columns[j]
			receiver.excel.SetCellValue("Sheet1", columnIndex+rowIndex, column)
		}
	}
}

func (receiver *Excel) SaveAs(path string) error {
	return receiver.excel.SaveAs(path)
}

func (receiver *Excel) WriteTo(writer io.Writer) (int64, error) {
	return receiver.excel.WriteTo(writer)
}

func (receiver *Excel) WriteToBuffer() (*bytes.Buffer, error) {
	return receiver.excel.WriteToBuffer()
}
