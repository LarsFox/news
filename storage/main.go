package main

import (
	"flag"
	"log"

	"github.com/LarsFox/news/storage/dbs"
	"github.com/LarsFox/news/storage/queues"
)

var (
	optDBPath = flag.String("db-path", "", "Путь до файла с БД")

	optNatsAddr    = flag.String("nats-addr", "", "Адрес nats.io")
	optGetNewsSubj = flag.String("get-news-subj", "", "Тема получения новостей")
)

func main() {
	flag.Parse()

	dbm, err := dbs.NewManager(*optDBPath)
	if err != nil {
		log.Fatal(err)
	}

	qm, err := queues.NewManager(*optNatsAddr, *optGetNewsSubj, dbm)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(qm.Listen())
}
