package main

import (
	"log"
	"parking_lot/parking"
	p "parking_lot/parking"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r parking.Repository
	p.ForeverSleep(2*time.Second, func(_ int) (err error) {
		//DebugDatabaseURL := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", "user", "pass", "localhost", "dbname")
		//DEBUG: pass DebugDatabaseURL rather cfg.DatabaseURL
		r, err = parking.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	s := parking.NewService(r)
	// connect service to server and start server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Fatal(parking.ListenGRPC(s, 5566))
		wg.Done()
	}()
	log.Println("GRPC Listening on port 5566...")
	wg.Add(1)
	go func() {
		log.Fatal(parking.ListenREST(s, 8080, 5566))
		wg.Done()
	}()
	log.Println("REST Listening on port 8080...")
	wg.Wait()
}
