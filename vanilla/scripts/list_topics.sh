#!/bin/zsh

if [[ "l" == $2 ]]; then
    docker exec -it $1 kafka-topics.sh --bootstrap-server=localhost:9092 --list
elif [[ "i" == $2 ]]; then
    docker exec -it $1 kafka-topics.sh --bootstrap-server=localhost:9092 --describe --topic $3
fi


