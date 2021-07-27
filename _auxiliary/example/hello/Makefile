.PHONY: proto data build

proto:
	for d in proto; do \
		for f in $$d/**/*.proto; do \
		    protoc  --roc_out=plugins=roc:. $$f; \
			echo compiled: $$f; \
		done; \
#	    for f in $$d/**/*.pb.go; do \
#    		    protoc-go-inject-tag -input=$$f;\
#    			echo compiled: $$f; \
#    	done; \
	done

deploy:
	./bin/build.sh; \
	./bin/deploy.sh; \
	./bin/run.sh

# build编译
b:
	./bin/build.sh


# 重启本地
r:
	./bin/build.sh; \
	./bin/stop.sh;\
	./bin/restart.sh;\

# 停止所有
s:
	./bin/stop.sh

# 本地连接远程，重启远程
rr:
	./bin/run.sh

restart:
	./bin/restart.sh



