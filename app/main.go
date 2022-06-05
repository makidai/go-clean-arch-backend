package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/makidai/go-clean-arch-backend/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/makidai/go-clean-arch-backend/article/delivery/http/middleware"
	_articleRepo "github.com/makidai/go-clean-arch-backend/article/repository/postgres"
	_articleUcase "github.com/makidai/go-clean-arch-backend/article/usecase"
	_authorRepo "github.com/makidai/go-clean-arch-backend/author/repository/postgres"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	dbConn, err := sql.Open("postgres", connection)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	arR := _articleRepo.NewPostgresArticleRepository(dbConn)
	auR := _authorRepo.NewPostgresAuthorRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(arR, auR, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
