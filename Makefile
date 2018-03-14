.PHONY: release

release:
	docker pull karalabe/xgo-latest
	go get -u -v github.com/karalabe/xgo
	mkdir -p bin 
	xgo -go "1.10" -dest bin ${LDFLAGS} -targets linux/amd64,linux/arm-6,darwin/amd64,windows/amd64 github.com/de0gee/basestation
	