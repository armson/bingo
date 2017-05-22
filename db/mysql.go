package db

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/fm"
    "fmt"
)

var Mysql Db = Db{}
type Tx struct{}

func (this *Db) Register(dsn string) {
    this.driver = "mysql"
    this.dsn = dsn

    var err error
    this.db, err = sql.Open(Mysql.driver, Mysql.dsn)
    if err != nil {
        panic(err.Error())
    }
}

func (this *Db) SetMaxIdleConns(maxIdleConns int64) {
    this.db.SetMaxIdleConns(int(maxIdleConns))
}

func (this *Db) SetMaxOpenConns(maxOpenConns int64) {
    this.db.SetMaxOpenConns(int(maxOpenConns))
}

func (this Db) String() string {
    return fmt.Sprintf("%s[%s][%+v]", this.driver, this.dsn, this.db)
}

func (this *Db) Query(sql string)  []map[string]string {
    //fmt.Println(sql)
    rows, err := this.db.Query(sql)
    if err != nil {
        panic(err.Error())
    }
    return this.queryFormat(rows)
}

func (this *Db) Fetch(sql string) (map[string]string) {
    rows := this.Query(sql)
    if len(rows) > 0 {
        return rows[0]
    }
    return nil
}

func (this *Db) Excute(sql string) (lastInsertId, afftectedRows int64) {
    fmt.Printf("%s\n\n", sql)
    res, err := this.db.Exec(sql)
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
func (this *Db) Prepare(sql string, args ...interface{}) (lastInsertId, afftectedRows int64){
    stm, err := this.db.Prepare(sql)
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

func (this *Db) Update(table string, data map[string]string, where string) (afftectedRows int64){
    sql := []string{"UPDATE ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    if where != "" {
        sql = append(sql, " WHERE ", where)
    }
    _, afftectedRows = this.Excute(fm.StringJoin(sql...))
    return
}
func (this *Db) Insert(table string, data map[string]string) (lastInsertId, afftectedRows int64){
    sql := []string{"INSERT INTO ",table," SET "}
    for k, v := range data {
        sql = append(sql, k,"='",v,"'"," , ")
    }
    sql = sql[:len(sql)-1]
    lastInsertId, afftectedRows = this.Excute(fm.StringJoin(sql...))
    return
}
//用法 
// tx := db.Mysql.Begin()
// for i := 0; i < 1000; i ++ {
//     tx.Exec("INSERT INTO user SET username = 'zhangfumu', age = 36")
// }
// tx.Commit()

func (this *Db) Begin() *sql.Tx {
    tx , err := this.db.Begin()
    if err != nil {
        panic(err.Error())
    }
    return tx
}

func (this *Db) queryFormat(rows *sql.Rows) (data []map[string]string) {
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



func (this *Db) Close() error {
    return this.db.Close()
}

func (this *Db) Db() *sql.DB {
    return this.db
}






