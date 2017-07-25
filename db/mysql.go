package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/armson/bingo/utils"
    "time"
)

type BinMysql struct{
    *sql.DB
}
var (
    Mysql *BinMysql
    MysqlGroup = map[string]*BinMysql{}
)

// 参考 go-sql-driver中dsn_test.go中的例子
func (_ *BinMysql) Register(group, dsn ,charset, timeout, readTimeout , writeTimeout string,
	maxIdleConn , maxOpenConn int , connMaxLifetime time.Duration ){
    dsn = utils.String.Join(
		dsn,
		"?charset=" ,charset,
		"&timeout=" ,timeout,
		"&readTimeout=" ,readTimeout,
		"&writeTimeout=" ,writeTimeout,
	)
	db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err.Error())
    }
	db.SetMaxIdleConns(maxIdleConn)
	db.SetMaxOpenConns(maxOpenConn)
	db.SetConnMaxLifetime(connMaxLifetime)

    if group == "db0" {
        Mysql = &BinMysql{db}
    }
    MysqlGroup[group] = &BinMysql{db}
}

func (_ *BinMysql) Use(group string) *BinMysql {
    if v , ok := MysqlGroup[group]; ok {
        return v
    }
    return MysqlGroup["db0"]
}

func (bin *BinMysql) Query(sql string)  []map[string]string {
	rows, err := bin.DB.Query(sql)
	if err != nil {
		panic(err.Error())
	}
	return QueryFormat(rows)
}

func (bin *BinMysql) Fetch(sql string) (map[string]string) {
	rows := bin.Query(sql+" Limit 1")
	if len(rows) > 0 {
		return rows[0]
	}
	return nil
}

func (bin *BinMysql) Execute(sql string) (lastInsertId, affectedRows int64) {
    res, err := bin.DB.Exec(sql)
    if err != nil {
        panic(err.Error())
    }
	affectedRows, err = res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }
    lastInsertId, err = res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    return
}

// Prepare的性能低于Execute
func (bin *BinMysql) Prepare(sql string, args ...interface{}) (lastInsertId, affectedRows int64){
    stm, err := bin.DB.Prepare(sql)
    if err != nil {
        panic(err.Error())
    }
    res, err := stm.Exec(args...)
    if err != nil {
        panic(err.Error())
    }
	affectedRows, err = res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }

    lastInsertId, err = res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    return
}

func (bin *BinMysql) Update(table string, data map[string]string, where string) (affectedRows int64){
    sql := []string{"UPDATE ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    if where != "" {
        sql = append(sql, " WHERE ", where)
    }
    _, affectedRows = bin.Execute(utils.String.Join(sql...))
    return
}
func (bin *BinMysql) Insert(table string, data map[string]string) (lastInsertId, affectedRows int64){
    sql := []string{"INSERT INTO ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    lastInsertId, affectedRows = bin.Execute(utils.String.Join(sql...))
    return
}
//用法
// tx := db.Mysql.Begin()
// for i := 0; i < 1000; i ++ {
//     tx.Exec("INSERT INTO user SET username = 'zhangfumu', age = 36")
// }
// tx.Commit()

func (bin *BinMysql) Begin() *sql.Tx {
    tx , err := bin.DB.Begin()
    if err != nil {
        panic(err.Error())
    }
    return tx
}

func QueryFormat(rows *sql.Rows) (data []map[string]string) {
    columns, _ := rows.Columns()
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        rows.Scan(scanArgs...)
        record := make(map[string]string)
        for i, col := range values {
            if col != nil {
                record[columns[i]] = string(col.([]byte))
            }
        }
        data = append(data, record)
    }
    return
}
func (bin *BinMysql) Close() error {
	return bin.DB.Close()
}

func (bin *BinMysql) Db() *sql.DB {
	return bin.DB
}











