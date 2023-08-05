package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

func run() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)

	rmq, err := initRabbitMQ()
	if err != nil {
		log.Fatalf("run: failed to init rabbitmq: %v", err)
	}
	defer rmq.Shutdown()

	// err = rmq.PublishWithDelay("user.event.publish", []byte("Haii guys!!"), 10000)
	// if err != nil {
	// 	log.Fatalf("run: failed to publish into rabbitmq: %v", err)
	// }

	for {
	}
}

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(os.Getenv("EXCHANGE")))
	})
	mux.HandleFunc("/schedule/request", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		if len(param) != 0 {
			host := param["host"][0]
			id := param["id"][0]
			gdate := param["date"][0]
			timestamp, err := strconv.ParseInt(param["timestamp"][0], 0, 64)

			data := map[string]interface{}{
				"id":   id,
				"host": host,
				"date": gdate,
			}

			jsonData, err := json.Marshal(data)

			rmq, err := initRabbitMQ()
			if err != nil {
				log.Fatalf("run: failed to init rabbitmq: %v", err)
			}
			defer rmq.Shutdown()

			err = rmq.PublishWithDelay("user.event.publish", []byte(jsonData), timestamp)
			if err != nil {
				log.Fatalf("run: failed to publish into rabbitmq: %v", err)
			}

			w.Write([]byte("success"))
		} else {
			w.Write([]byte("failed"))
		}

	})

	server := http.Server{
		Addr:    ":5790",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
	go runServer()
	run()

}
