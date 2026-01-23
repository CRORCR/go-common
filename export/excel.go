package export

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter Excel导出器
type ExcelExporter struct {
	file *excelize.File
}

// NewExcelExporter 创建Excel导出器
func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{
		file: excelize.NewFile(),
	}
}

// ExportToExcel 导出数据到Excel
func (e *ExcelExporter) ExportToExcel(data []map[string]interface{}, columns []string, sheetName string) error {
	// 创建或获取工作表
	index, err := e.file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 准备表头和数据，批量写入
	if len(data) > 0 {
		// 构建完整的行数据（包括表头）
		allRows := make([][]interface{}, len(data)+1)
		// 表头行
		headerRow := make([]interface{}, len(columns))
		for i, col := range columns {
			headerRow[i] = col
		}
		allRows[0] = headerRow

		// 数据行
		for rowIdx, row := range data {
			rowData := make([]interface{}, len(columns))
			for colIdx, col := range columns {
				value := row[col]
				if value == nil {
					value = ""
				}
				rowData[colIdx] = value
			}
			allRows[rowIdx+1] = rowData
		}

		// 批量写入数据到工作表
		for i, rowData := range allRows {
			cell := fmt.Sprintf("A%d", i+1)
			if err = e.file.SetSheetRow(sheetName, cell, &rowData); err != nil {
				return err
			}
		}
	} else {
		// 如果没有数据，只设置表头
		headerRow := make([]interface{}, len(columns))
		for i, col := range columns {
			headerRow[i] = col
		}
		if err = e.file.SetSheetRow(sheetName, "A1", &headerRow); err != nil {
			return err
		}
	}

	// 设置表头样式
	headerStyle, err := e.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	// 应用表头样式
	endCol := columnIndexToLetter(len(columns) - 1)
	if err = e.file.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", endCol), headerStyle); err != nil {
		return err
	}

	// 自动调整列宽
	for i := range columns {
		colLetter := columnIndexToLetter(i)
		if err = e.file.SetColWidth(sheetName, colLetter, colLetter, 15); err != nil {
			return err
		}
	}

	// 设置活动工作表
	e.file.SetActiveSheet(index)

	return nil
}

// SaveToFile 保存到文件
func (e *ExcelExporter) SaveToFile(filename string) error {
	// 删除默认的Sheet1
	if e.file.GetSheetName(0) == "Sheet1" {
		e.file.DeleteSheet("Sheet1")
	}
	return e.file.SaveAs(filename)
}

// SaveToBuffer 保存到缓冲区
func (e *ExcelExporter) SaveToBuffer() ([]byte, error) {
	// 删除默认的Sheet1
	if e.file.GetSheetName(0) == "Sheet1" {
		e.file.DeleteSheet("Sheet1")
	}
	buffer, err := e.file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// columnIndexToLetter 将列索引转换为Excel列字母
func columnIndexToLetter(index int) string {
	var result strings.Builder
	index++ // Excel列从1开始

	for index > 0 {
		index--
		result.WriteByte(byte('A' + index%26))
		index /= 26
	}

	// 反转字符串
	runes := []rune(result.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}
