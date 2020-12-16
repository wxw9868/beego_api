package sqlinit

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB
var err error

func initDatabase() {
	connStr := "postgresql://bank_user:dfc.1465@104.168.163.18:26257/bank?sslmode=require"
	//connStr := "postgresql://root@104.168.163.18:26257/bank?ssl=true&sslmode=require&sslrootcert=/databases/certs/ca.crt&sslkey=/databases/certs/client.root.key&sslcert=/databases/certs/client.root.crt"
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	//defer db.Close()
}
