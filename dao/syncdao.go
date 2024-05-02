package dao

import (
	"database/sql"
	"xmserver/com"
)

func InsertOrUpdateSync(s *com.Sync) error {
	insertOrUpdateSQL := "INSERT INTO synch (table, row_id, t) VALUES (?, ?, ?) ON CONFLICT(table,row_id) DO UPDATE SET t = excluded.t"
	_, err := db.Exec(insertOrUpdateSQL, s.Table, s.RowID, s.T)
	return err
}

func InsertSync(s *com.Sync) (sql.Result, error) {
	insertSQL := `INSERT INTO synch (table, row_id, t) VALUES (?, ?, ?)`
	result, err := db.Exec(insertSQL, s.Table, s.RowID, s.T)
	return result, err
}

func GetSync(id int64) (*com.Sync, error) {
	s := &com.Sync{}
	querySQL := `SELECT id, table, row_id, t FROM synch WHERE id = ?`
	row := db.QueryRow(querySQL, id)
	err := row.Scan(&s.ID, &s.Table, &s.RowID, &s.T)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func UpdateSync(s *com.Sync) (sql.Result, error) {
	updateSQL := `UPDATE synch SET table = ?, row_id = ?, t = ? WHERE id = ?`
	result, err := db.Exec(updateSQL, s.Table, s.RowID, s.T, s.ID)
	return result, err
}

func DeleteSync(id int64) (sql.Result, error) {
	deleteSQL := `DELETE FROM synch WHERE id = ?`
	result, err := db.Exec(deleteSQL, id)
	return result, err
}
