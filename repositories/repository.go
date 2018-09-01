package repositories

import (
	"fmt"
	"os"

	"github.com/go-xorm/xorm"

	_ "github.com/go-sql-driver/mysql"
)

var engine *xorm.Engine

func init() {
	fmt.Println("init xorm Engine")
	var err error

	var hostname string
	var dbname string
	var username string
	var password string
	var port string

	// TODO
	// LOCAL / PRODUCTION 分離したい
	hostname = os.Getenv("LOCAL_DB_HOSTNAME")
	dbname = os.Getenv("LOCAL_DB_DBNAME")
	username = os.Getenv("LOCAL_DB_USERNAME")
	password = os.Getenv("LOCAL_DB_PASSWORD")
	port = os.Getenv("LOCAL_DB_PORT")

	engine, err = xorm.NewEngine("mysql", username+":"+password+"@tcp("+hostname+":"+port+")/"+dbname)

	if err != nil {
		panic(err)
	}
}
