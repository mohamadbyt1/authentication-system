package main
import (
	"log"
	"github.com/mohamadbyt1/authentication-system/api"
	"github.com/mohamadbyt1/authentication-system/storage"
)
func main() {
	db, err := storage.NewDb()
	if err != nil {
		log.Fatal(err)
	}
	addr := "0.0.0.0:8080"
	s := api.NewApiServer(addr, db)
	serverErr := s.Start()
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
