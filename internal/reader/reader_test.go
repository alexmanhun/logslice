package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestReadFile_ValidJSON(t *testing.T) {
	p := writeTempFile(t, `{"level":"info","msg":"started"}
{"level":"error","msg":"failed"}
`)

	recs, errs := ReadFile(p)
	var got []Record
	for r := range recs {
		got = append(got, r)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 records, got %d", len(got))
	}
	if got[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", got[0].Fields["level"])
	}
	if got[1].Line != 2 {
		t.Errorf("expected line 2, got %d", got[1].Line)
	}
}

func TestReadFile_SkipsNonJSON(t *testing.T) {
	p := writeTempFile(t, "not json\n{\"level\":\"warn\"}\nalso not json\n")

	recs, errs := ReadFile(p)
	var got []Record
	for r := range recs {
		got = append(got, r)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 record, got %d", len(got))
	}
}

func TestReadFile_MissingFile(t *testing.T) {
	_, errs := ReadFile("/nonexistent/path/file.log")
	if err := <-errs; err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestReadFiles_MultipleFiles(t *testing.T) {
	p1 := writeTempFile(t, "{\"src\":\"a\"}\n")
	p2 := writeTempFile(t, "{\"src\":\"b\"}\n{\"src\":\"c\"}\n")

	recs, errs := ReadFiles([]string{p1, p2})
	var got []Record
	for r := range recs {
		got = append(got, r)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 records, got %d", len(got))
	}
}
