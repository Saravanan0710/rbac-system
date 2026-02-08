package main

import (
    "database/sql"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"

    _ "modernc.org/sqlite"
)

func backupDB(src string) (string, error) {
    ts := time.Now().Format("20060102T150405")
    dst := fmt.Sprintf("%s.bak-%s", src, ts)
    in, err := os.Open(src)
    if err != nil {
        return "", err
    }
    defer in.Close()
    out, err := os.Create(dst)
    if err != nil {
        return "", err
    }
    defer out.Close()
    if _, err := io.Copy(out, in); err != nil {
        return "", err
    }
    return dst, nil
}

func main() {
    // assume running from rbac-backend root so rbac.db is present
    dbPath := "rbac.db"
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        log.Fatalf("database file %s not found in cwd", dbPath)
    }

    bak, err := backupDB(dbPath)
    if err != nil {
        log.Fatalf("backup failed: %v", err)
    }
    log.Printf("backup created: %s", bak)

    db, err := sql.Open("sqlite", "file:rbac.db?_foreign_keys=on")
    if err != nil {
        log.Fatalf("open db: %v", err)
    }
    defer db.Close()

    tx, err := db.Begin()
    if err != nil {
        log.Fatalf("begin tx: %v", err)
    }

    // Delete tasks first (tasks reference projects)
    stmts := []string{
        "DELETE FROM tasks;",
        "DELETE FROM project_assignments;",
        "DELETE FROM projects;",
        "DELETE FROM users WHERE role != 'ADMIN';",
    }

    for _, s := range stmts {
        if _, err := tx.Exec(s); err != nil {
            tx.Rollback()
            log.Fatalf("exec failed (%s): %v", s, err)
        }
    }

    if err := tx.Commit(); err != nil {
        log.Fatalf("commit failed: %v", err)
    }

    abs, _ := filepath.Abs(dbPath)
    log.Printf("cleanup completed against %s", abs)
}
