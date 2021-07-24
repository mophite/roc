#!/usr/bin/env bash

kill -9 $(ps -ef|grep "\./api" |awk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')

kill -9 $(ps -ef|grep "\./srv" |awk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')