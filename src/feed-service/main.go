package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/feed-service/api/router"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/src/feed-service/feed"
)

func main() {
	port := os.Getenv("PORT")
	feedController := feed.NewDefaultController()
	handler := router.New(feedController)
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
