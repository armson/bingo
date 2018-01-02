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
var (
	mysqlCluster = map[string]*sql.DB{}
	mysqlDb = []string{} //clusterID单独存储，主要是为了方便指定默认的数据库
	mysqlAlias = make(map[string]string) //{"cache1":"db0", "cache2":"db1", "cache3":"db2","db0":"db0", "db1":"db1", "db2":"db2"}
	isValid = false
)


func init()  {
	mysqlDb = config.Slice("mysql","dbs");
	if  len(mysqlDb) == 0 { return }

	charset := config.String("mysql","charset")
	timeout := config.String("mysql","timeout")
	readTimeout := config.String("mysql","readTimeout")
	writeTimeout := config.String("mysql","writeTimeout")
	maxIdleConns := config.Int("mysql","maxIdleConns")
	maxOpenConns := config.Int("mysql","maxOpenConns")
	connMaxLifetime := config.Time("mysql","connMaxLifetime")

	for _ , id := range mysqlDb {
		register(
			id,
			config.String("mysql:"+id, "dsn"),
			charset, timeout, readTimeout, writeTimeout, maxIdleConns, maxOpenConns, connMaxLifetime,
		)

		alias := config.String("mysql:"+id, "alias")
		if alias == "" { alias = id }
		mysqlAlias[alias] = id
		mysqlAlias[id] = id
	}

	// mysql可用
	isValid = true
}

type binMysql struct {
	tracer bingo.Tracer
	id string
}
// 参考 go-sql-driver中dsn_test.go中的例子
func register(group, dsn ,charset, timeout, readTimeout , writeTimeout string,
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
	mysqlCluster[group] = db
}

func New(trance bingo.Tracer) *binMysql {
	return &binMysql{
		tracer:	trance,
		id:		mysqlDb[0],
	}
}

func Valid() bool { return isValid }

func (my *binMysql) Use(name string) *binMysql {
	if id , ok := mysqlAlias[name]; ok {
		my.id = id
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
	return mysqlCluster[bin.id]
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

















