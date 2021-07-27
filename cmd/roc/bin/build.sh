#!/usr/bin/env bash

dir=`pwd`
ldflags="-X 'main.version=$(git log --pretty=format:'%H' -n 1)'"

build() {
	for d in $(ls ./$1); do
		echo "building $1/$d"

    if [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
        	pushd $dir/$1/$d > $dir/cmd/NUL
        	CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "$ldflags -w -s"
    else
      		pushd $dir/$1/$d >/dev/null
      		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$ldflags -w -s"
    fi

		if [ ! -d $dir/cmd/$d/ ]; then
          mkdir -p $dir/cmd/$d/
    fi

    if [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
          mv $dir/$1/$d/$d.exe $dir/cmd/$d/
    else
      	  rsync  $dir/$1/$d/$d $dir/cmd/$d/
    fi

    if [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
          popd > $dir/cmd/NUL
          rm $dir/cmd/NUL
    else
      		popd >/dev/null
    fi
	done

}

build app/api
build app/srv