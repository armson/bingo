package mysql

import (
	"github.com/armson/bingo"
	"github.com/armson/bingo/config"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/armson/bingo/utils"
	"encoding/json"
    "time"
)

type binMysql struct {
	tracer bingo.Tracer
	id string
}
var mysqlGroup = map[string]*sql.DB{}


// 参考 go-sql-driver中dsn_test.go中的例子
func Register(group, dsn ,charset, timeout, readTimeout , writeTimeout string,
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
	mysqlGroup[group] = db
}

func New(trance bingo.Tracer) *binMysql {
	return &binMysql{
		tracer:	trance,
		id:		"db0",
	}
}

func (my *binMysql) Use(group string) *binMysql {
	if _ , ok := mysqlGroup[group]; ok {
		my.id = group
	}
	return my
}

func (bin *binMysql) Logs(cost time.Duration, sql string, result string) {
	if config.Bool("default","enableLog") && config.Bool("mysql","enableLog") {
		bin.tracer.Logs("Mysql", utils.String.Join(
			"Cost: ", cost.String(),
			" Sql: ", sql,
			" Result: ", result,
		),
		)
	}
}

func (bin *binMysql) Query(sql string) ([]map[string]string) {
	begin := time.Now()
	rows, err := bin.pool().Query(sql)
	if err != nil {
		bin.Logs(time.Since(begin), sql, err.Error())
		panic(err.Error())
	}
	result := queryFormat(rows)
	s, _ := json.Marshal(result)
	bin.Logs(time.Since(begin), sql, string(s))
	return result
}

func (bin *binMysql) Fetch(sql string) (map[string]string) {
	rows := bin.Query(sql+" Limit 1")
	if len(rows) > 0 { return rows[0] }
	return nil
}

func (bin *binMysql) Execute(sql string) (lastInsertId, affectedRows int64) {
	begin := time.Now()
	res, err := bin.pool().Exec(sql)
	if err != nil {
		bin.Logs(time.Since(begin), sql, err.Error())
		panic(err.Error())
	}
	affectedRows, err = res.RowsAffected()
	if err != nil {
		bin.Logs(time.Since(begin), sql, err.Error())
		panic(err.Error())
	}
	lastInsertId, err = res.LastInsertId()
	if err != nil {
		bin.Logs(time.Since(begin), sql, err.Error())
		panic(err.Error())
	}
	result := utils.String.Join(
		utils.Int.String(affectedRows), " ", utils.Int.String(lastInsertId),
	)
	bin.Logs(time.Since(begin), sql, result)
	return
}

func (bin *binMysql) Update(table string, data map[string]string, where string) (affectedRows int64){
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

func (bin *binMysql) Insert(table string, data map[string]string) (lastInsertId, affectedRows int64){
	sql := []string{"INSERT INTO ",table," SET "}
	for k, v := range data {
		sql = append(sql, k,"='",v,"'"," , ")
	}
	sql = sql[:len(sql)-1]
	lastInsertId, affectedRows = bin.Execute(utils.String.Join(sql...))
	return
}

func (bin *binMysql) pool() *sql.DB {
	return mysqlGroup[bin.id]
}

func (bin *binMysql) Close() error {
	return bin.pool().Close()
}

func (bin *binMysql) Db() *sql.DB {
	return bin.pool()
}

func queryFormat(rows *sql.Rows) (data []map[string]string) {
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




//
//type BinMysql struct{
//	*sql.DB
//}
//func (bin *BinMysql) Query(sql string)  []map[string]string {
//	rows, err := bin.DB.Query(sql)
//	if err != nil {
//		panic(err.Error())
//	}
//	return QueryFormat(rows)
//}
//
//func (bin *BinMysql) Fetch(sql string) (map[string]string) {
//	rows := bin.Query(sql+" Limit 1")
//	if len(rows) > 0 {
//		return rows[0]
//	}
//	return nil
//}
//
//func (bin *BinMysql) Execute(sql string) (lastInsertId, affectedRows int64) {
//    res, err := bin.DB.Exec(sql)
//    if err != nil {
//        panic(err.Error())
//    }
//	affectedRows, err = res.RowsAffected()
//    if err != nil {
//        panic(err.Error())
//    }
//    lastInsertId, err = res.LastInsertId()
//    if err != nil {
//        panic(err.Error())
//    }
//    return
//}
//
//// Prepare的性能低于Execute
//func (bin *BinMysql) Prepare(sql string, args ...interface{}) (lastInsertId, affectedRows int64){
//    stm, err := bin.DB.Prepare(sql)
//    if err != nil {
//        panic(err.Error())
//    }
//    res, err := stm.Exec(args...)
//    if err != nil {
//        panic(err.Error())
//    }
//	affectedRows, err = res.RowsAffected()
//    if err != nil {
//        panic(err.Error())
//    }
//
//    lastInsertId, err = res.LastInsertId()
//    if err != nil {
//        panic(err.Error())
//    }
//    return
//}
//
//func (bin *BinMysql) Update(table string, data map[string]string, where string) (affectedRows int64){
//    sql := []string{"UPDATE ",table," SET "}
//    for k, v := range data {
//        sql = append(sql, k,"='",v,"'"," , ")
//    }
//    sql = sql[:len(sql)-1]
//    if where != "" {
//        sql = append(sql, " WHERE ", where)
//    }
//    _, affectedRows = bin.Execute(utils.String.Join(sql...))
//    return
//}
//func (bin *BinMysql) Insert(table string, data map[string]string) (lastInsertId, affectedRows int64){
//    sql := []string{"INSERT INTO ",table," SET "}
//    for k, v := range data {
//        sql = append(sql, k,"='",v,"'"," , ")
//    }
//    sql = sql[:len(sql)-1]
//    lastInsertId, affectedRows = bin.Execute(utils.String.Join(sql...))
//    return
//}
////用法
//// tx := db.Mysql.Begin()
//// for i := 0; i < 1000; i ++ {
////     tx.Exec("INSERT INTO user SET username = 'zhangfumu', age = 36")
//// }
//// tx.Commit()
//
//func (bin *BinMysql) Begin() *sql.Tx {
//    tx , err := bin.DB.Begin()
//    if err != nil {
//        panic(err.Error())
//    }
//    return tx
//}
//













