default:
	go run -race ./*.go \
		--db-path=.tmp/seed.db \
		--nats-addr=:5555 \
		--get-news-subj=gns

test:
	go test -timeout 30s github.com/LarsFox/news/storage/dbs
