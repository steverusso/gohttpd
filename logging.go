// This is free and unencumbered software released
// into the public domain. Please see the UNLICENSE
// file or unlicense.org for more information.
package gohttpd

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Logger interface {
	Log(r *http.Request)
	Close() error
}

func NewLogger(f string) (Logger, error) {
	if strings.HasSuffix(f, ".sqlite") {
		return newSQLiteLogger(f)
	} else if f != "" {
		return newFileLogger(f)
	}
	return nil, nil
}

type fileLogger struct {
	io.WriteCloser
}

func newFileLogger(fpath string) (Logger, error) {
	if fpath == "stdout" {
		return &fileLogger{os.Stdout}, nil
	}

	f, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &fileLogger{f}, nil
}

func (l *fileLogger) Log(r *http.Request) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(l, "%s :: %s :: %s\n", now, r.Host, r.URL)
}

func (l *fileLogger) Close() error {
	if l.WriteCloser == os.Stdout {
		return nil
	}
	return l.WriteCloser.Close()
}

type sqliteLogger struct {
	*sql.DB
}

func newSQLiteLogger(fpath string) (Logger, error) {
	const initQuery = `
		CREATE TABLE IF NOT EXISTS page_views (
			id       INTEGER PRIMARY KEY,
			host     TEXT,
			resource TEXT,
			reqtime  TEXT DEFAULT CURRENT_TIMESTAMP
		)`

	db, err := sql.Open("sqlite3", fpath)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	if _, err := db.Exec(initQuery); err != nil {
		return nil, fmt.Errorf("executing init query: %w", err)
	}

	return &sqliteLogger{db}, nil
}

func (l *sqliteLogger) Log(r *http.Request) {
	const logQuery = `
		INSERT INTO page_views (host, resource)
		VALUES (?, ?)`

	if _, err := l.Exec(logQuery, r.Host, r.URL.String()); err != nil {
		log.Printf("error inserting page view: %v", err)
	}
}
