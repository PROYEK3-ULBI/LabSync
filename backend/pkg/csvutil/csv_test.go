package csvutil

import (
	"strings"
	"testing"
)

func TestSanitizeCell(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Normal Text", "Normal Text"},
		{"=cmd|' /C calc'!A0", "'=cmd|' /C calc'!A0"},
		{"+12345", "'+12345"},
		{"-SUM(A1:A10)", "'-SUM(A1:A10)"},
		{"@SUM(A1:A10)", "'@SUM(A1:A10)"},
		{"\tTabPrefix", "'\tTabPrefix"},
		{"  =FormulaWithSpace", "'  =FormulaWithSpace"},
		{"", ""},
	}

	for _, tt := range tests {
		got := SanitizeCell(tt.input)
		if got != tt.expected {
			t.Errorf("SanitizeCell(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestGenerateCSV(t *testing.T) {
	headers := []string{"ID", "Nama", "Keterangan"}
	rows := [][]string{
		{"1", "Budi", "Hadir"},
		{"2", "=malicious_formula()", "Selesai"},
	}

	csvBytes, err := GenerateCSV(headers, rows)
	if err != nil {
		t.Fatalf("GenerateCSV error: %v", err)
	}

	csvStr := string(csvBytes)
	if !strings.Contains(csvStr, "ID,Nama,Keterangan") {
		t.Errorf("expected headers in output")
	}
	if !strings.Contains(csvStr, "'=malicious_formula()") {
		t.Errorf("expected sanitized formula in output, got: %s", csvStr)
	}
}
