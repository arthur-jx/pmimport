export PATH := $(GOPATH)/bin:$(PATH)
#LDFLAGS := -s -w
LDFLAGS := -w

all: build

build: app

app:
#       env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o pmimport_darwin_amd64 .
#       env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o pmimport_linux_386 .
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o pmimport_linux_amd64 .
#       env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags "$(LDFLAGS)" -o pmimport_linux_arm .
#       env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 cc=arm-fsl-linux-gnueabi-gcc go build -ldflags "$(LDFLAGS)" -o pmimport_arm-v5 ./
#       env CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabihf-gcc go build -ldflags "$(LDFLAGS)" -o pmimport_arm-v7l ./
#       env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o pmimport_linux_arm64 .
#       env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "$(LDFLAGS)" -o pmimport_windows_386.exe .
#       env CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "$(LDFLAGS)" -o pmimport_windows_386.exe .
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o pmimport_windows_amd64.exe .

