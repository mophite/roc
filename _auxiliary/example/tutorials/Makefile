.PHONY: proto data build

proto:
	for d in proto; do \
		for f in $$d/**/*.proto; do \
		    protoc  --roc_out=plugins=roc:. $$f; \
			echo compiled: $$f; \
		done; \
	done

# build
build:
	./bin/build.sh


# build and stop ,then restart
run:
	./bin/build.sh; \
	./bin/stop.sh;\
	./bin/restart.sh;\

# stop
stop:
	./bin/stop.sh

# restart
restart:
	./bin/restart.sh