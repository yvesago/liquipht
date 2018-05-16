VERSION := 0.0.2

LDFLAGS := -ldflags "-s"
#GOARCH = 386

BINSERVER := liquipht
TWMINER := twminer

ARCHS := 386 amd64
archs = $(word 1, $@)

.PHONY: all
all: all32 all64 shasum

.PHONY: all32 all64
all32: 386 windows linux darwin zip

all64: amd64 windows linux darwin

.PHONY: $(ARCHS)
$(ARCHS): 
	@echo $(archs)	
	$(eval GOARCH = $(archs))


.PHONY: windows
windows:
	mkdir -p release
	GOOS=windows GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(BINSERVER)-v$(VERSION)-win-$(GOARCH).exe server.go
	GOOS=windows GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(TWMINER)-v$(VERSION)-win-$(GOARCH).exe twitter.go

.PHONY: linux
linux:
	mkdir -p release
	GOOS=linux GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(BINSERVER)-v$(VERSION)-linux-$(GOARCH) server.go
	GOOS=linux GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(TWMINER)-v$(VERSION)-linux-$(GOARCH) twitter.go

.PHONY: darwin
darwin:
	mkdir -p release
	GOOS=darwin GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(BINSERVER)-v$(VERSION)-darwin-$(GOARCH) server.go
	GOOS=darwin GOARCH=$(GOARCH) go build $(LDFLAGS) -o release/$(TWMINER)-v$(VERSION)-darwin-$(GOARCH) twitter.go

.PHONY: zip
zip:
	mkdir -p release
	zip -q -r release/views-v$(VERSION).zip views/

.PHONY: shasum
shasum:
	cd release; sha1sum *-v$(VERSION)* > sha1sum-$(VERSION).txt
