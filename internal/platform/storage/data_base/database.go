package data_base

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

var db *sql.DB

func Connect() (*sql.DB, error) {
	//err := os.Remove("./db.sqlite")
	//if err != nil {
	//	return nil, err
	//}

	db, _ = sql.Open("sqlite3", "./db.sqlite")
	db.SetMaxOpenConns(2)
	db.SetConnMaxLifetime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Start() (*sql.DB, error) {
	sqlStmt := `CREATE TABLE IF NOT EXISTS requests (
    id INTEGER PRIMARY KEY, 
    ip VARCHAR(50) default NULL, 
    header TEXT, 
    command TEXT, 
    status TEXT, 
    created_at TEXT default CURRENT_TIMESTAMP, 
    updated_at TEXT default CURRENT_TIMESTAMP);
	
	CREATE TABLE IF NOT EXISTS commands (
	id INTEGER PRIMARY KEY,
	"whitelist"	TEXT,
	"blacklist"	TEXT,
	"created_at" TEXT default CURRENT_TIMESTAMP,
	"updated_at" TEXT default CURRENT_TIMESTAMP);

	INSERT INTO commands (whitelist, blacklist) VALUES ('', '"sudo",
			"shutdown", "whoami", "cd ", "reboot", "systemctl", "rm ", "rmdir ", "mkdir ", "touch ",
			"mv ", "cp ", "cat ", "less", "more", "head", "tail", "find", "grep", "awk", "sed",
			"sort ", "uniq", "wc", "diff", "patch", "tar ", "gzip ", "gunzip", "bzip2", "unzip", "zip ",
			"chown ", "chmod ", "chgrp ", "chattr ", "chcon ", "chroot ", "chvt", "chsh", "chfn", "chage ",
			"chpasswd", "vi", "vim", "nano", "emacs", "ged", "gedit", "kate", "kwrite", "kedit", "python ",
			"perl ", "ruby ", "php ", "java ", "javac ", "gcc ", "g++ ", "make ", "cmake ", "clang ", "clang++ ", "rustc ",
			"go ", "node ", "npm ", "yarn ", "pip ", "pip3 ", "pipenv ", "docker ", "docker-compose ", "docker-machine ",
			"docker-swarm ", "docker-credential ", "dockerd ", "docker-init ", "docker-proxy ", "docker-runc ", "dockerd ",
			"/.", "sh ", "ssh ", "scp ", "sftp ", "rsync ", "curl ", "wget ", "aria2c ", "aria2 ", "aria2c ", "aria2 "');
    `

	res, err := db.Exec(sqlStmt)
	if err != nil {
		return db, err
	}

	log.Printf("Result: %v", res)

	return db, nil
}

func Close() error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}

// Insert - insert data to database
func Insert(ctx *gin.Context, status int) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO requests (header, command, ip, status) values(?, ?, ?, ?)")
	if err != nil {
		return false, err
	}
	log.Printf("Header: %v", ctx.Request.Header)

	_, err = stmt.Exec(ctx.Request.Header.Get("User-Agent"), ctx.Query("cmd"), ctx.ClientIP(), status)
	if err != nil {
		return false, err
	}

	err = stmt.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetAll() (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM requests")
	if err != nil {
		return nil, err
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func GetBlackList() (*sql.Rows, error) {
	rows, err := db.Query("SELECT blacklist FROM commands LIMIT 1")
	if err != nil {
		return nil, err
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	return rows, nil
}
