.PHONY: all
all: windows linux

.PHONY: linux-i386
linux-i386:
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o ./out/winreg-tasks-linux-i386 -trimpath ./cmd

.PHONY: windows-i386
windows-i386:
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o ./out/winreg-tasks-windows-i386.exe -trimpath ./cmd

.PHONY: windows-amd64
windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./out/winreg-tasks-windows-amd64.exe -trimpath ./cmd

.PHONY: linux-amd64
linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out/winreg-tasks-linux-amd64 -trimpath ./cmd

.PHONY: windows
windows: windows-i386 windows-amd64

.PHONY: linux
linux: linux-i386 linux-amd64

.PHONY: i386
i386: linux-i386 windows-i386

.PHONY: amd64
amd64: linux-amd64 windows-amd64

.PHONY: targets
targets:
	@LC_ALL=C $(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/(^|\n)# Files(\n|$$)/,/(^|\n)# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | grep -E -v -e '^[^[:alnum:]]' -e '^$@$$'
