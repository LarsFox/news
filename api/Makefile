default:
	go run -race ./*.go \
		--api-addr=:5545 \
		--nats-addr=:5555 \
		--get-news-subj=gns

g200:
	curl localhost:5545/news/how-to-seed/

g404:
	curl localhost:5545/news/article-404/
