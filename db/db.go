package db

import (
    "database/sql"
)

type Db struct{
    driver string
    dsn    string
    db     *sql.DB
}

type dber interface{
    Register(dsn string)
}