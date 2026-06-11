package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRotatingWriter_CreatesDateFile(t *testing.T) {
	dir := t.TempDir()
	fixed := time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC)
	w := &rotatingWriter{dir: dir, maxAge: 0, nowFunc: func() time.Time { return fixed }}

	_, err := w.Write([]byte("hello\n"))
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	path := filepath.Join(dir, "2024-06-10.log")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("file not created: %v", err)
	}
	if string(data) != "hello\n" {
		t.Fatalf("expected 'hello\\n', got %q", string(data))
	}
}

func TestRotatingWriter_RotatesOnDateChange(t *testing.T) {
	dir := t.TempDir()
	yesterday := time.Date(2024, 1, 14, 23, 59, 59, 0, time.UTC)
	today := time.Date(2024, 1, 15, 0, 0, 1, 0, time.UTC)
	now := yesterday
	w := &rotatingWriter{dir: dir, maxAge: 0, nowFunc: func() time.Time { return now }}

	if _, err := w.Write([]byte("day1\n")); err != nil {
		t.Fatalf("Write day1 error: %v", err)
	}
	now = today
	if _, err := w.Write([]byte("day2\n")); err != nil {
		t.Fatalf("Write day2 error: %v", err)
	}

	file1 := filepath.Join(dir, "2024-01-14.log")
	file2 := filepath.Join(dir, "2024-01-15.log")
	if _, err := os.Stat(file1); err != nil {
		t.Fatalf("file1 not found: %v", err)
	}
	if _, err := os.Stat(file2); err != nil {
		t.Fatalf("file2 not found: %v", err)
	}
	data1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("read file1 error: %v", err)
	}
	data2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("read file2 error: %v", err)
	}
	if string(data1) != "day1\n" {
		t.Fatalf("file1 content: %q", string(data1))
	}
	if string(data2) != "day2\n" {
		t.Fatalf("file2 content: %q", string(data2))
	}
}

func TestRotatingWriter_CleanupOldFiles(t *testing.T) {
	dir := t.TempDir()
	// 创建 3 个旧文件（超过 maxAge=7 天）
	for _, name := range []string{"2024-01-01.log", "2024-01-02.log", "2024-01-03.log"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("old"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// 创建 1 个不应被删除的文件（在 maxAge 范围内）
	recent := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	if err := os.WriteFile(filepath.Join(dir, recent+".log"), []byte("recent"), 0644); err != nil {
		t.Fatal(err)
	}

	w := &rotatingWriter{dir: dir, maxAge: 7, nowFunc: time.Now}
	if _, err := w.Write([]byte("today\n")); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	for _, name := range []string{"2024-01-01.log", "2024-01-02.log", "2024-01-03.log"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Errorf("expected %s to be deleted", name)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, recent+".log")); err != nil {
		t.Errorf("expected recent file to be kept: %v", err)
	}
}

func TestRotatingWriter_MaxAgeZeroSkipsCleanup(t *testing.T) {
	dir := t.TempDir()
	old := "2020-01-01.log"
	if err := os.WriteFile(filepath.Join(dir, old), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	w := &rotatingWriter{dir: dir, maxAge: 0, nowFunc: time.Now}
	if _, err := w.Write([]byte("x")); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, old)); err != nil {
		t.Errorf("maxAge=0 should not delete files, but %s was deleted", old)
	}
}

func TestInit_WritesDateFile(t *testing.T) {
	dir := t.TempDir()
	Init(Config{WriteFile: true, LogDir: dir, MaxAge: 7})
	t.Cleanup(func() { Init(Config{}) })
	Info("test message")
	today := time.Now().Format("2006-01-02")
	path := filepath.Join(dir, today+".log")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("log file not created: %v", err)
	}
	if !strings.Contains(string(data), "test message") {
		t.Fatalf("log file does not contain 'test message', got: %s", string(data))
	}
}
