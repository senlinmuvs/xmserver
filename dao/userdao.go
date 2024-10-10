package dao

import (
	"encoding/json"
	"fmt"
	"log"
	"xmserver/com"
)

func UpsertRecord(user, table, k, v, jsonData string) {
	udb, ok := userDBS[user]
	if !ok {
		panic(com.Err1("UpsertRecord user db err", user))
	}
	// 将JSON数据解析到map中，以便构建SQL语句
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Println("UpsertRecord unmarshalling err", jsonData)
		panic(com.Err1("UpsertRecord unmarshalling JSON err", err))
	}

	// 构建SQL语句的字段和值部分
	fields := ""
	values := ""
	updateStmt := ""
	var args []interface{}
	for field, value := range data {
		if fields != "" {
			fields += ", "
			values += ", "
			updateStmt += ", "
		}
		fields += field
		values += "?"
		updateStmt += fmt.Sprintf("%s = ?", field)
		args = append(args, value)
	}

	// 检查记录是否存在
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, k)
	if err := udb.QueryRow(checkQuery, v).Scan(&exists); err != nil {
		fmt.Println("UpsertRecord error sql", checkQuery, user, table, k, v, "args", len(args), args)
		panic(com.Err1("UpsertRecord error checking if record exists", err))
	}

	if exists {
		updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", table, updateStmt, k)
		args = append(args, v)
		_, err := udb.Exec(updateQuery, args...)
		if err != nil {
			fmt.Println("UpsertRecord error sql", updateQuery, user, table, k, v, "args", len(args), args)
			panic(com.Err1("UpsertRecord error executing update query", err))
		}
	} else {
		insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, fields, values)
		_, err := udb.Exec(insertQuery, args...)
		if err != nil {
			fmt.Println("UpsertRecord error sql", insertQuery, user, table, k, v, "args", len(args), args)
			panic(com.Err1("UpsertRecord error executing insert query", err))
		}
	}
}

func DelRecord(user, table, k, v string) {
	udb, ok := userDBS[user]
	if !ok {
		panic(com.Err1("DelRecord user db err", user))
	}
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, k)
	stmt, err := udb.Prepare(sql)
	if err != nil {
		log.Println("DelRecord user", user, "table", table, k, v, "err1", err, sql)
		panic(com.Err1("DelRecord err1", err))
	}
	defer stmt.Close()

	_, err = stmt.Exec(v)
	if err != nil {
		log.Println("DelRecord user", user, "table", table, k, v, "err2", err, sql)
		panic(com.Err1("DelRecord err2", err))
	}
}
func ExeSql(user, sql string) {
	udb, ok := userDBS[user]
	if !ok {
		panic(com.Err1("ExeSql user db err", user))
	}
	_, err := udb.Exec(sql)
	if err != nil {
		log.Println("ExcSql user", user, "err", err, "sql", sql)
		panic(com.Err1("ExeSql err", err))
	}
}
