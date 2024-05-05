package dao

import (
	"database/sql"
	"log"
	"os"
	"xmserver/com"
)

var (
	db      *sql.DB
	userDBS = map[string]*sql.DB{}
)

func InitDao(sysDbFile string) {
	db = connectDatabase(com.Cfg.XMS.DataDir + "/" + sysDbFile)
	fileInfos, err := os.ReadDir(com.Cfg.XMS.DataDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			u := fileInfo.Name()
			userDbFile := com.Cfg.XMS.DataDir + "/" + u + "/xxmoon.data"
			d, e := sql.Open("sqlite3", userDbFile)
			if e != nil {
				log.Fatal(e, "db open fail1", userDbFile)
			}
			e = d.Ping()
			if e != nil {
				log.Fatal(e, "db open fail2", userDbFile)
			}
			userDBS[u] = d
			log.Println("open db suc", userDbFile)
		}
	}
	log.Println("init dao done user dir size", len(userDBS))
}

func connectDatabase(dbfile string) *sql.DB {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS sync (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tab TEXT,
		rowid Integer,
		t Integer
	);
	CREATE TABLE IF NOT EXISTS env (
		k TEXT PRIMARY KEY,
		v TEXT
	);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
func Close() {
	db.Close()
}
