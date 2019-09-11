default:

nats:
	docker run -p 5555:4444 nats -p 4444

proto:
	mkdir -p storage/vendor/github.com/LarsFox && cp -r newsgrpc storage/vendor/github.com/LarsFox
	mkdir -p web/vendor/github.com/LarsFox && cp -r newsgrpc web/vendor/github.com/LarsFox
