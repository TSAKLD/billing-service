package main

import (
	"avitoTZ/api"
	"avitoTZ/bootstrap"
	"avitoTZ/repository"
	"avitoTZ/service"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	c, err := bootstrap.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := repository.DBConnect(c)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.New(db)

	serv := service.NewUserService(repo)

	r := api.NewServer(serv, c.HTTPPort)

	err = r.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
