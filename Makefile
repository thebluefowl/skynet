build:
	env GOOS=linux GOARCH=386 go build -o dist/skynet main.go
	env GOOS=linux GOARCH=386 go build -o dist/seed seed/main.go
	env GOOS=linux GOARCH=386 go build -o dist/migrate migration/main.go

deploy:
	scp dist/skynet vishnu@157.245.97.106:/home/vishnu/skynet/app
	scp dist/seed vishnu@157.245.97.106:/home/vishnu/skynet/seed
	scp dist/migrate vishnu@157.245.97.106:/home/vishnu/skynet/migrate
	scp dist/IP2LOCATION-LITE-DB5.IPV6.BIN vishnu@157.245.97.106:/home/vishnu/skynet/IP2LOCATION-LITE-DB5.IPV6.BIN