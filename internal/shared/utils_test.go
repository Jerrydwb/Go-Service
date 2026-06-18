package shared

import (
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "0.00"},
		{1234.567, "1234.57"},
		{-50.5, "-50.50"},
		{0.999, "1.00"},
	}

	for _, tt := range tests {
		result := FormatNumber(tt.input)
		if result != tt.expected {
			t.Errorf("FormatNumber(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetFloatValue(t *testing.T) {
	val := 42.5
	tests := []struct {
		input    *float64
		expected float64
	}{
		{nil, 0},
		{&val, 42.5},
	}

	for _, tt := range tests {
		result := GetFloatValue(tt.input)
		if result != tt.expected {
			t.Errorf("GetFloatValue(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestFormatPeriodo(t *testing.T) {
	tests := []struct {
		dateType string
		dateStart string
		expected string
	}{
		{"", "", "Periodo no especificado"},
		{"anual", "2024", "Periodo Anual: 2024"},
		{"mes", "05/2024", "Periodo Mensual: Mayo de 2024"},
		{"mensual", "01/2023", "Periodo Mensual: Enero de 2023"},
		{"otro", "algo", "Periodo: algo"},
	}

	for _, tt := range tests {
		result := FormatPeriodo(tt.dateType, tt.dateStart)
		if result != tt.expected {
			t.Errorf("FormatPeriodo(%q, %q) = %q, want %q", tt.dateType, tt.dateStart, result, tt.expected)
		}
	}
}

func TestFlexibleStringUnmarshalJSON(t *testing.T) {
	tests := []struct {
		json     string
		expected string
	}{
		{`"hello"`, "hello"},
		{`123`, "123"},
		{`12.0`, "12.0"}, // json.Number atrapa primero → preserva formato
	}

	for _, tt := range tests {
		var fs FlexibleString
		if err := fs.UnmarshalJSON([]byte(tt.json)); err != nil {
			t.Errorf("UnmarshalJSON(%s) error: %v", tt.json, err)
			continue
		}
		if string(fs) != tt.expected {
			t.Errorf("UnmarshalJSON(%s) = %q, want %q", tt.json, string(fs), tt.expected)
		}
	}
}

func TestToLatin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"ascii", "hello", "hello"},
		{"tildes", "Denominación", "Denominaci\xf3n"},
		{"enie", "Año", "A\xf1o"},
		{"multiple", "Comercio & Distribución S.A.C.", "Comercio & Distribuci\xf3n S.A.C."},
		{"pregunta", "¿Cómo?", "\xbfC\xf3mo?"},
		{"exclamacion", "¡Atención!", "\xa1Atenci\xf3n!"},
		{"u_umlaut", "Pingüino", "Ping\xfcino"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLatin(tt.input)
			if result != tt.expected {
				t.Errorf("ToLatin(%q) = %q (len=%d), want %q (len=%d)",
					tt.input, result, len(result), tt.expected, len(tt.expected))
			}
		})
	}
}
