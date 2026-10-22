package csvutil

import (
	"bytes"
	"encoding/csv"
	"strings"
)

// SanitizeCell mencegah serangan Formula Injection (CSV Injection)
// Menetralkan nilai yang diawali tanda =, +, -, @, \t, atau \r dengan menambahkan prefix single quote (')
func SanitizeCell(val string) string {
	if val == "" {
		return ""
	}
	trimmed := strings.TrimLeft(val, " ")
	if len(trimmed) > 0 {
		firstChar := trimmed[0]
		if firstChar == '=' || firstChar == '+' || firstChar == '-' || firstChar == '@' || firstChar == '\t' || firstChar == '\r' {
			return "'" + val
		}
	}
	return val
}

// GenerateCSV membuat byte array CSV dari headers dan rows dengan sanitasi formula injection
func GenerateCSV(headers []string, rows [][]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)

	sanitizedHeaders := make([]string, len(headers))
	for i, h := range headers {
		sanitizedHeaders[i] = SanitizeCell(h)
	}
	if err := writer.Write(sanitizedHeaders); err != nil {
		return nil, err
	}

	for _, row := range rows {
		sanitizedRow := make([]string, len(row))
		for i, cell := range row {
			sanitizedRow[i] = SanitizeCell(cell)
		}
		if err := writer.Write(sanitizedRow); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
