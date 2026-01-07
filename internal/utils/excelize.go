package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// 修正 excelize v2.4.0 要求 time.Time 必须是 UTC 时区的限制
// https://github.com/360EntSecGroup-Skylar/excelize/issues/409
func FixExcelizeTime(t time.Time) time.Time {
	t, _ = time.ParseInLocation("2006-01-02 15:04:05", t.Local().Format("2006-01-02 15:04:05"), time.UTC)
	return t
}

type ExcelRowReader struct {
	// 列标题 => 列下标（0开始）
	ColumnIndexes map[string]int
}

func (r ExcelRowReader) ColumnIndex(columnTitles ...string) (int, bool) {
	for _, columnTitle := range columnTitles {
		if columnIndex, ok := r.ColumnIndexes[columnTitle]; ok {
			return columnIndex, true
		}
	}
	return 0, false
}

func (r ExcelRowReader) ColumnNumber(columnTitles ...string) int {
	if col, ok := r.ColumnIndex(columnTitles...); ok {
		return col + 1
	}
	return 0
}

// 根据预先解析的列下标，读取某行对应列的值（去左右空格）
// 如果列不在标题行中，或当前行对应列未填，返回空字符串
func (r ExcelRowReader) Read(row []string, columnTitles ...string) string {
	for _, columnTitle := range columnTitles {
		if columnIndex, ok := r.ColumnIndexes[columnTitle]; ok && columnIndex < len(row) {
			return strings.TrimSpace(row[columnIndex])
		}
	}
	return ""
}

func (r ExcelRowReader) ReadUppercase(row []string, columnTitles ...string) string {
	return strings.ToUpper(r.Read(row, columnTitles...))
}

func (r ExcelRowReader) ReadLowercase(row []string, columnTitles ...string) string {
	return strings.ToLower(r.Read(row, columnTitles...))
}

func (r ExcelRowReader) PrefixColumnIndex(prefixes ...string) (int, bool) {
	for _, prefix := range prefixes {
		for columnTitle, columnIndex := range r.ColumnIndexes {
			if strings.HasPrefix(columnTitle, prefix) {
				return columnIndex, true
			}
		}
	}
	return 0, false
}

func (r ExcelRowReader) PrefixColumnNumber(prefixes ...string) int {
	if col, ok := r.PrefixColumnIndex(prefixes...); ok {
		return col + 1
	}
	return 0
}

func (r ExcelRowReader) PrefixRead(row []string, prefixes ...string) string {
	for _, prefix := range prefixes {
		for columnTitle, columnIndex := range r.ColumnIndexes {
			if strings.HasPrefix(columnTitle, prefix) && columnIndex < len(row) {
				return strings.TrimSpace(row[columnIndex])
			}
		}
	}
	return ""
}

func (r ExcelRowReader) PrefixReadUppercase(row []string, prefixes ...string) string {
	return strings.ToUpper(r.PrefixRead(row, prefixes...))
}

func (r ExcelRowReader) PrefixReadLowercase(row []string, prefixes ...string) string {
	return strings.ToLower(r.PrefixRead(row, prefixes...))
}

// 解析标题行，此行每个非空值（去左右空格后）都会用于指示一个字段的列下标，用户读取后续行中的值
// 例如：["标识", "名称\n", " "] => {"标识": 0, "名称": 1}
//
//goland:noinspection GoUnusedExportedFunction
func NewExcelRowReader(titleLine []string) ExcelRowReader {
	columnIndexes := make(map[string]int, len(titleLine))
	for columnIndex, columnTitle := range titleLine {
		columnTitle = strings.TrimSpace(columnTitle)
		if columnTitle != "" {
			columnIndexes[columnTitle] = columnIndex
		}
	}
	return ExcelRowReader{ColumnIndexes: columnIndexes}
}

type GenerateTemplateExcelColumnDef struct {
	Title    string
	Type     string // string, int, float
	DropList []string
}

// 按列定义创建 sheet，并写入标题行
func GenerateTemplateExcel(workbook *excelize.File, sheetName string, columnDefs []GenerateTemplateExcelColumnDef) (err error) {
	_, err = workbook.NewSheet(sheetName)
	if err != nil {
		return
	}

	if len(columnDefs) == 0 {
		return
	}

	titleRow := make([]any, 0, len(columnDefs))
	for _, columnDef := range columnDefs {
		titleRow = append(titleRow, columnDef.Title)
	}

	if err = workbook.SetSheetRow(sheetName, "A1", &titleRow); err != nil {
		return
	}

	return
}

// 添加单元格样式和数据校验
func GenerateTemplateExcelStyles(workbook *excelize.File, sheetName string, columnDefs []GenerateTemplateExcelColumnDef, rowCount int) (err error) {
	if len(columnDefs) == 0 {
		return
	}

	stringColumns := make([]string, 0, len(columnDefs))
	validations := make([]*excelize.DataValidation, 0, len(columnDefs))

	maxRowNumber := strconv.Itoa(rowCount)

	for columnIdx, columnDef := range columnDefs {
		if len(columnDef.DropList) > 0 {
			columnName, err := excelize.ColumnNumberToName(columnIdx + 1)
			if err != nil {
				return err
			}

			dv := excelize.NewDataValidation(true)
			dv.Sqref = columnName + "2:" + columnName + maxRowNumber
			dv.SetDropList(columnDef.DropList)
			validations = append(validations, dv)
		} else if columnDef.Type != "" {
			columnName, err := excelize.ColumnNumberToName(columnIdx + 1)
			if err != nil {
				return err
			}
			switch columnDef.Type {
			case "string":
				stringColumns = append(stringColumns, columnName)
			}
		}
	}

	if len(stringColumns) > 0 {
		var style int
		style, err = workbook.NewStyle(&excelize.Style{NumFmt: 49})
		if err != nil {
			return
		}
		for _, columnName := range stringColumns {
			err = workbook.SetCellStyle(sheetName, columnName+"2", columnName+maxRowNumber, style)
			if err != nil {
				return
			}
		}
	}

	for _, dv := range validations {
		if err = workbook.AddDataValidation(sheetName, dv); err != nil {
			return err
		}
	}

	return
}

func RemoveTemplateExcelTables(workbook *excelize.File, sheetName string) error {
	tables, err := workbook.GetTables(sheetName)
	if err != nil {
		return err
	}
	for _, table := range tables {
		if err := workbook.DeleteTable(table.Name); err != nil {
			return err
		}
	}
	return nil
}

// 按预定义样式在 A1 到 (col, row) 创建一个 table，其名称和 sheetName 相同
// 一般来说:
// 导出数据时在 A1 到 (columnCount, dataCount + 1) 创建一个 table（"+ 1" 因为标题需要占一行）
// 导出模板时在 A1 到 (columnCount, 2) 创建一个 table（"2" 因为 table 需要至少两行）
func CreateTemplateExcelTable(workbook *excelize.File, sheetName string, col, row int) error {
	return CreateTemplateExcelTableWithName(workbook, sheetName, sheetName, col, row)
}

func CreateTemplateExcelTableWithName(workbook *excelize.File, sheetName string, tableName string, col, row int) error {
	cell, err := excelize.CoordinatesToCellName(col, row)
	if err != nil {
		return err
	}
	showHeaderRow := true
	showRowStripes := false
	return workbook.AddTable(sheetName, &excelize.Table{
		Range:          "A1:" + cell,
		Name:           tableName,
		StyleName:      "TableStyleLight13",
		ShowHeaderRow:  &showHeaderRow,
		ShowRowStripes: &showRowStripes,
	})
}
