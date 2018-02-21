all:
	./scripts/build.sh

install:
	./scripts/install.sh
	
test: 
	cd bin/ && \
	cat ../tests/test1.txt | ./master
	#./../scripts/test.sh
