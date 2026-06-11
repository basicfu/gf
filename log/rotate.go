package log

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type rotatingWriter struct {
	dir     string
	maxAge  int
	date    string
	file    *os.File
	mu      sync.Mutex
	nowFunc func() time.Time
}

func (w *rotatingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.nowFunc == nil {
		w.nowFunc = time.Now
	}
	today := w.nowFunc().Format("2006-01-02")
	if today != w.date {
		if err := w.rotate(today); err != nil {
			return 0, err
		}
	}
	return w.file.Write(p)
}

func (w *rotatingWriter) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}

func (w *rotatingWriter) rotate(today string) error {
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	if err := os.MkdirAll(w.dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(w.dir, today+".log")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w.file = f
	w.date = today
	if w.maxAge > 0 {
		now, _ := time.Parse("2006-01-02", today)
		w.cleanup(now)
	}
	return nil
}

func (w *rotatingWriter) cleanup(now time.Time) {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	cutoff := now.AddDate(0, 0, -w.maxAge)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if filepath.Ext(name) != ".log" {
			continue
		}
		dateStr := name[:len(name)-4]
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			_ = os.Remove(filepath.Join(w.dir, name))
		}
	}
}
