package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jayvib/clean-architecture/config"
	"github.com/sirupsen/logrus"
	"net/url"
)

func New(conf *config.Config) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		conf.Database.User, conf.Database.Pass,
		conf.Database.Host, conf.Database.Port,
		conf.Database.Name,
	)
	logrus.Println(conn)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Manila")
	dsn := fmt.Sprintf("%s?%s", conn, val.Encode())
	dbConn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}
