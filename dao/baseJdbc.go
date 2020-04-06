package dao

import (
	"bvtracker/g"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var SqlDB *sql.DB

//golang这边实现的连接池只提供了SetMaxOpenConns和SetMaxIdleConns方法进行连接池方面的配置。
// 在使用的过程中有一个问题就是数据库本身对连接有一个超时时间的设置，如果超时时间到了数据库会单方面断掉连接，此时再用连接池内的连接进行访问就会出错。
func Init() {
	//SqlDB, _ = sql.Open("mysql", "root:Bv123456@tcp(47.92.68.38:3306)/bvtracker?charset=utf8")
	SqlDB, _ = sql.Open("mysql", g.Config().MysqlCon)

	//SqlDB, _ = sql.Open("mysql", "root:123456@tcp(localhost:3306)/bvtracker?charset=utf8")
	////用于设置最大打开的连接数，默认值为0表示不限制
	//SqlDB.SetMaxOpenConns(10)
	////用于设置闲置的连接数。
	//SqlDB.SetMaxIdleConns(10)
	SqlDB.SetConnMaxLifetime(100 * time.Second)
	SqlDB.Ping()
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//select
func SelectMysqlData(sql string) ([]map[string]interface{}, error) {
	rows, err := SqlDB.Query(sql)
	checkErr(err)
	defer rows.Close()

	columns, _ := rows.Columns()
	if err != nil {
		return nil, err
	}
	tableData := make([]map[string]interface{}, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

//insert
func InsertMysqlData(sql string) (int64, bool) {

	res, err := SqlDB.Exec(sql)
	checkErr(err)
	id, err := res.LastInsertId()
	if err != nil {
		return 0, false
	} else {
		return id, true
	}
}

//update
func UpdateMysqlData(sql string) (int64, bool) {

	res, err := SqlDB.Exec(sql)
	checkErr(err)
	id, err := res.RowsAffected()
	if err != nil {
		return 0, false
	} else {
		return id, true
	}
}

//delete
func DeleteMysqlData(sql string) bool {

	res, err := SqlDB.Exec(sql)
	checkErr(err)
	_, err = res.RowsAffected()
	if err != nil {
		return false
	} else {
		return true
	}
}

// count
func CountRecord(sql string) int {
	num := 0
	rows, err := SqlDB.Query(sql)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		num++
	}
	return num
}
