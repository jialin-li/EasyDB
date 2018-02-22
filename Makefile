Tags = DEBUG

all:
	./scripts/build.sh $(Tags)

install:
	./scripts/install.sh
	
test: 
	./scripts/test.sh

release:
	./scripts/build.sh
.PHONY: release all test install
