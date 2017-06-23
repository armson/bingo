package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/armson/bingo/utils"
    "time"
)


type binMysql struct{
    *sql.DB
}
var (
    Mysql *binMysql
    MysqlGroup = map[string]*binMysql{}
)

// 参考 go-sql-driver中dsn_test.go中的例子
func (this *binMysql) Register(group, dsn string, params map[string]string){
    dsn = utils.String.Join(dsn, "?" , utils.Map.BuildQuery(params))
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        panic(err.Error())
    }
    if group == "db0" {
        Mysql = &binMysql{db}
    }
    MysqlGroup[group] = &binMysql{db}
}

func (this *binMysql) Use(group string) *binMysql {
    if v , ok := MysqlGroup[group]; ok {
        return v
    }
    return MysqlGroup["db0"]
}

func (this *binMysql) SetMaxIdleConns(maxIdleConns int) {
    this.DB.SetMaxIdleConns(maxIdleConns)
}

func (this *binMysql) SetMaxOpenConns(maxOpenConns int) {
    this.DB.SetMaxOpenConns(maxOpenConns)
}

func (this *binMysql) SetConnMaxLifetime(t time.Duration) {
    this.DB.SetConnMaxLifetime(t)
}

func (this *binMysql) Query(sql string)  []map[string]string {
    rows, err := this.DB.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    return this.queryFormat(rows)
}

func (this *binMysql) Fetch(sql string) (map[string]string) {
    rows := this.Query(sql)
    if len(rows) > 0 {
        return rows[0]
    }
    return nil
}

func (this *binMysql) Excute(sql string) (lastInsertId, afftectedRows int64) {
    res, err := this.DB.Exec(sql)
    if err != nil {
        panic(err.Error())
    }
    afftectedRows, err = res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }
    lastInsertId, err = res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    return
}

// Prepare的性能低于Excute
func (this *binMysql) Prepare(sql string, args ...interface{}) (lastInsertId, afftectedRows int64){
    stm, err := this.DB.Prepare(sql)
    if err != nil {
        panic(err.Error())
    }
    res, err := stm.Exec(args...)
    if err != nil {
        panic(err.Error())
    }
    afftectedRows, err = res.RowsAffected()
    if err != nil {
        panic(err.Error())
    }

    lastInsertId, err = res.LastInsertId()
    if err != nil {
        panic(err.Error())
    }
    return
}

func (this *binMysql) Update(table string, data map[string]string, where string) (afftectedRows int64){
    sql := []string{"UPDATE ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    if where != "" {
        sql = append(sql, " WHERE ", where)
    }
    _, afftectedRows = this.Excute(utils.String.Join(sql...))
    return
}
func (this *binMysql) Insert(table string, data map[string]string) (lastInsertId, afftectedRows int64){
    sql := []string{"INSERT INTO ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    lastInsertId, afftectedRows = this.Excute(utils.String.Join(sql...))
    return
}
//用法 
// tx := db.Mysql.Begin()
// for i := 0; i < 1000; i ++ {
//     tx.Exec("INSERT INTO user SET username = 'zhangfumu', age = 36")
// }
// tx.Commit()

func (this *binMysql) Begin() *sql.Tx {
    tx , err := this.DB.Begin()
    if err != nil {
        panic(err.Error())
    }
    return tx
}

func (this *binMysql) queryFormat(rows *sql.Rows) (data []map[string]string) {
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

func (this *binMysql) Close() error {
    return this.DB.Close()
}

func (this *binMysql) Db() *sql.DB {
    return this.DB
}










