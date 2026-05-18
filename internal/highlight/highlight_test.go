package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestLevelColor_KnownLevels(t *testing.T) {
	cases := []struct {
		level    string
		wantCode string
	}{
		{"error", highlight.Red},
		{"fatal", highlight.Red},
		{"warn", highlight.Yellow},
		{"warning", highlight.Yellow},
		{"info", highlight.Green},
		{"debug", highlight.Cyan},
		{"trace", highlight.Cyan},
		{"unknown", highlight.Gray},
	}
	for _, tc := range cases {
		t.Run(tc.level, func(t *testing.T) {
			got := highlight.LevelColor(tc.level)
			if got != tc.wantCode {
				t.Errorf("LevelColor(%q) = %q, want %q", tc.level, got, tc.wantCode)
			}
		})
	}
}

func TestLevelColor_CaseInsensitive(t *testing.T) {
	if highlight.LevelColor("ERROR") != highlight.Red {
		t.Error("expected ERROR to map to Red")
	}
	if highlight.LevelColor("Info") != highlight.Green {
		t.Error("expected Info to map to Green")
	}
}

func TestColorizeLevel_ContainsUpperCase(t *testing.T) {
	result := highlight.ColorizeLevel("info")
	if !strings.Contains(result, "INFO") {
		t.Errorf("ColorizeLevel(\"info\") should contain \"INFO\", got %q", result)
	}
	if !strings.Contains(result, highlight.Reset) {
		t.Error("ColorizeLevel should contain reset code")
	}
}

func TestColorizeKey_Bold(t *testing.T) {
	result := highlight.ColorizeKey("level")
	if !strings.Contains(result, highlight.Bold) {
		t.Error("ColorizeKey should contain bold code")
	}
	if !strings.Contains(result, "level") {
		t.Error("ColorizeKey should contain the key name")
	}
}

func TestStrip_RemovesAnsi(t *testing.T) {
	input := highlight.ColorizeLevel("error")
	stripped := highlight.Strip(input)
	if strings.Contains(stripped, "\033") {
		t.Errorf("Strip should remove escape codes, got %q", stripped)
	}
	if stripped != "ERROR" {
		t.Errorf("Strip result = %q, want \"ERROR\"", stripped)
	}
}

func TestStrip_PlainString(t *testing.T) {
	input := "no color here"
	if got := highlight.Strip(input); got != input {
		t.Errorf("Strip(%q) = %q, want unchanged", input, got)
	}
}
