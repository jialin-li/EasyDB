ARGS = -tags DEBUG

all:
	./scripts/build.sh $(ARGS)

install:
	./scripts/install.sh
	
# tests with debug statements
test: 
	make
	./scripts/make_test.sh

# tests without debug statements
rtest:
	make release
	./scripts/release_test.sh

release:
	./scripts/build.sh

clean:
	./scripts/clean.sh

.PHONY: release all test install clean
