package main

import (
	"flag"
	"log"

	"github.com/LarsFox/news/api/api"
	"github.com/LarsFox/news/api/queues"
)

var (
	optApiAddr = flag.String("api-addr", "", "Адрес сервиса")

	optNatsAddr    = flag.String("nats-addr", "", "Адрес nats.io")
	optGetNewsSubj = flag.String("get-news-subj", "", "Тема получения новостей")
)

func main() {
	flag.Parse()

	queuesM, err := queues.NewManager(*optNatsAddr, *optGetNewsSubj)
	if err != nil {
		log.Fatal(err)
	}

	apiM := api.NewManager(queuesM)

	log.Fatal(apiM.Listen(*optApiAddr))
}
