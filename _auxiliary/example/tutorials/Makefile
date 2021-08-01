.PHONY: proto data build

proto:
	for d in proto; do \
		for f in $$d/**/*.proto; do \
		    protoc  --roc_out=plugins=roc:. $$f; \
			echo compiled: $$f; \
		done; \
	done

# build
b:
	./bin/build.sh


# build and stop ,then restart
r:
	./bin/build.sh; \
	./bin/stop.sh;\
	./bin/restart.sh;\

# stop
s:
	./bin/stop.sh

# restart
rs:
	./bin/restart.sh



