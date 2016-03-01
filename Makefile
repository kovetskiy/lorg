test:
	go test -v

cover:
	go test -coverprofile cover 2>/dev/null
	go tool cover -html=cover -o cover.html

clean:
	rm cover cover.html
