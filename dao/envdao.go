package dao

import (
	"database/sql"
	"xmserver/com"
)

func InsertEnv(k, v string) {
	insertSQL := "INSERT INTO env (k, v) VALUES (?, ?)"
	_, err := db.Exec(insertSQL, k, v)
	if err != nil {
		panic(com.Err1(err.Error()))
	}
}

func InsertOrUpdateEnv(k, v string) {
	// 注意：这适用于支持ON CONFLICT语法的SQLite版本（3.24.0及以上）
	insertOrUpdateSQL := "INSERT INTO env (k, v) VALUES (?, ?) ON CONFLICT(k) DO UPDATE SET v = excluded.v"
	_, err := db.Exec(insertOrUpdateSQL, k, v)
	if err != nil {
		panic(com.Err1(err.Error()))
	}
}

func GetEnv(k string) string {
	var v string
	querySQL := `SELECT v FROM env WHERE k = ?`
	row := db.QueryRow(querySQL, k)
	err := row.Scan(&v)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		}
		panic(com.Err1(err.Error()))
	}
	return v
}

func UpdateEnv(k, v string) {
	updateSQL := `UPDATE env SET v = ? WHERE k = ?`
	_, err := db.Exec(updateSQL, v, k)
	if err != nil {
		panic(com.Err1(err.Error()))
	}
}

func DeleteEnv(k string) {
	deleteSQL := `DELETE FROM env WHERE k = ?`
	_, err := db.Exec(deleteSQL, k)
	if err != nil {
		panic(com.Err1(err.Error()))
	}
}
