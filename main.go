package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jayvib/app/config"
	"github.com/jayvib/app/internal/app/search/elasticsearch"
	jlog "github.com/jayvib/app/log"
	"github.com/olivere/elastic/v7"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	articlehttp "github.com/jayvib/app/article/delivery/http"
	articlerepo "github.com/jayvib/app/article/repository/mysql"
	articleusecase "github.com/jayvib/app/article/usecase"
	authorrepo "github.com/jayvib/app/author/repository/mysql"
	jwtmiddleware "github.com/jayvib/app/middleware/jwt"
	userhttp "github.com/jayvib/app/user/delivery/http"
	userrepo "github.com/jayvib/app/user/repository/mysql"
	usersearches "github.com/jayvib/app/user/search/elasticsearch"
	articlesearches "github.com/jayvib/app/article/search/elasticsearch"
	userusecase "github.com/jayvib/app/user/usecase"
	"github.com/sirupsen/logrus"
)

var conf *config.Config

var (
	Version     string
	Environment string
)

var (
	db *sql.DB
	es *elastic.Client
)

var (
	// ESIndexMappingFilename consist of index name as a key and
	// the file path to its equivalent mapping.
	ESIndexMappingFilename = map[string]string{
		"user":    "user.json",
		"article": "article.json",
	}
)

func init() {
	logrus.Println("Initializing...")
	var err error
	conf, err = config.New()
	if err != nil {
		panic(err)
	}
	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.AddHook(jlog.NewDebugHook())
		logrus.Info("Server running on DEBUG mode.")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	Environment = os.Getenv(config.AppEnvironmentKey)
}

func main() {
	printDBValues()
	printInfo()

	// ##########THIRD PARTY##########
	var err error
	db, err = newDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	es = newESClient()
	createESIndex()

	// ###########ROUTER##############
	e := gin.Default()
	api := e.Group("/")           // authentication not required
	authapi := e.Group("/api/v1") // authentication is required
	authapi.Use(jwtmiddleware.Authenticate(conf.JWTToken))

	// ##########Author###########
	authorRepo := authorrepo.New(db)

	// ###########User############
	userRepo := userrepo.New(db)
	userSearchEngine := usersearches.New(es)
	userUsecase := userusecase.New(userRepo, authorRepo, userSearchEngine)

	// ##########Article###########
	articleRepo := articlerepo.New(db)
	articleSE := articlesearches.New(es)
	articleUsecase := articleusecase.New(articleRepo, authorRepo, articleSE, time.Second*2)

	// Register handlers groups
	userhttp.RegisterHandlers(conf, api, userUsecase)
	articlehttp.RegisterHandlers(conf, authapi, articleUsecase)
	log.Fatal(e.Run(conf.Server.Address))
}

func printDBValues() {
	logrus.WithFields(logrus.Fields{
		"DBHOST": conf.Database.Host,
		"DBPORT": conf.Database.Port,
		"DBUSER": conf.Database.User,
		"DBPASS": conf.Database.Pass,
		"DBNAME": conf.Database.Name,
	}).Debug("Database Environment Values")
}

func printInfo() {
	fmt.Println("==========INFO===========")
	fmt.Println("Environment:", Environment)
}

func newDBConnection() (*sql.DB, error) {
	conn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		conf.Database.User, conf.Database.Pass,
		conf.Database.Host, conf.Database.Port,
		conf.Database.Name,
	)
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

func newESClient() *elastic.Client {
	esClient, err := elastic.NewClient(
		elastic.SetURL(conf.Elasticsearch.Servers...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	return esClient
}

func createESIndex() {
	for index, mappingFile := range ESIndexMappingFilename {
		err := elasticsearch.CreateIndex(es, index, mappingFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
