package export

import (
	"database/sql"
	"fmt"
	"strings"
)

// SQLExporter SQL导出器
type SQLExporter struct {
	builder strings.Builder
}

// NewSQLExporter 创建SQL导出器
func NewSQLExporter() *SQLExporter {
	return &SQLExporter{}
}

// ExportType SQL导出类型
type ExportType string

const (
	ExportInsert     ExportType = "insert"      // INSERT语句
	ExportUpdate     ExportType = "update"      // UPDATE语句
	ExportInsertOnly ExportType = "insert_only" // 仅INSERT数据
	ExportComplete   ExportType = "complete"    // 完整备份（含CREATE TABLE）
)

// ExportToSQL 导出数据为SQL
func (e *SQLExporter) ExportToSQL(
	tableName string,
	data []map[string]interface{},
	columns []string,
	exportType ExportType,
	tableSchema string,
	whereColumns []string,
) (string, error) {
	e.builder.Reset()

	// 添加注释
	e.builder.WriteString(fmt.Sprintf("-- SQL Export for table: %s\n", tableName))
	e.builder.WriteString(fmt.Sprintf("-- Export type: %s\n", exportType))
	e.builder.WriteString(fmt.Sprintf("-- Total rows: %d\n\n", len(data)))

	// 根据导出类型生成SQL
	switch exportType {
	case ExportComplete:
		// 完整备份：包含建表语句
		if tableSchema != "" {
			e.builder.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", tableName))
			e.builder.WriteString(tableSchema)
			e.builder.WriteString("\n\n")
		}
		fallthrough
	case ExportInsert, ExportInsertOnly:
		return e.generateInsertSQL(tableName, data, columns), nil
	case ExportUpdate:
		if len(whereColumns) == 0 {
			return "", fmt.Errorf("WHERE columns are required for UPDATE export")
		}
		return e.generateUpdateSQL(tableName, data, columns, whereColumns), nil
	default:
		return "", fmt.Errorf("unsupported export type: %s", exportType)
	}
}

// generateInsertSQL 生成INSERT语句
func (e *SQLExporter) generateInsertSQL(tableName string, data []map[string]interface{}, columns []string) string {
	if len(data) == 0 {
		return e.builder.String()
	}

	// 生成列名部分
	columnsPart := "(" + strings.Join(wrapBackticks(columns), ", ") + ")"

	// 批量插入，每100条数据一组
	batchSize := 100
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end]
		e.builder.WriteString(fmt.Sprintf("INSERT INTO `%s` %s VALUES\n", tableName, columnsPart))

		for j, row := range batch {
			values := make([]string, len(columns))
			for k, col := range columns {
				values[k] = formatValue(row[col])
			}

			e.builder.WriteString("(")
			e.builder.WriteString(strings.Join(values, ", "))
			e.builder.WriteString(")")

			if j < len(batch)-1 {
				e.builder.WriteString(",\n")
			} else {
				e.builder.WriteString(";\n\n")
			}
		}
	}

	return e.builder.String()
}

// generateUpdateSQL 生成UPDATE语句
// whereColumns: 用作WHERE条件的列名数组，支持多个字段组合
func (e *SQLExporter) generateUpdateSQL(tableName string, data []map[string]interface{}, columns []string, whereColumns []string) string {
	// 将whereColumns转为map以便快速查找
	whereColMap := make(map[string]bool)
	for _, col := range whereColumns {
		whereColMap[col] = true
	}

	for _, row := range data {
		setParts := make([]string, 0)
		whereParts := make([]string, 0)

		for _, col := range columns {
			if whereColMap[col] {
				// WHERE条件字段
				whereParts = append(whereParts, fmt.Sprintf("`%s` = %s", col, formatValue(row[col])))
			} else {
				// SET更新字段
				setParts = append(setParts, fmt.Sprintf("`%s` = %s", col, formatValue(row[col])))
			}
		}

		if len(setParts) == 0 || len(whereParts) == 0 {
			continue
		}

		e.builder.WriteString(fmt.Sprintf("UPDATE `%s` SET ", tableName))
		e.builder.WriteString(strings.Join(setParts, ", "))
		e.builder.WriteString(" WHERE ")
		e.builder.WriteString(strings.Join(whereParts, " AND "))
		e.builder.WriteString(";\n")
	}

	return e.builder.String()
}

// formatValue 格式化值为SQL字符串
func formatValue(value interface{}) string {
	if value == nil {
		return "NULL"
	}

	switch v := value.(type) {
	case string:
		// 转义特殊字符
		escaped := strings.ReplaceAll(v, "\\", "\\\\")
		escaped = strings.ReplaceAll(escaped, "'", "\\'")
		escaped = strings.ReplaceAll(escaped, "\n", "\\n")
		escaped = strings.ReplaceAll(escaped, "\r", "\\r")
		return fmt.Sprintf("'%s'", escaped)
	case []byte:
		return fmt.Sprintf("'%s'", string(v))
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case sql.NullString:
		if v.Valid {
			return formatValue(v.String)
		}
		return "NULL"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// wrapBackticks 为列名添加反引号
func wrapBackticks(columns []string) []string {
	result := make([]string, len(columns))
	for i, col := range columns {
		result[i] = fmt.Sprintf("`%s`", col)
	}
	return result
}
