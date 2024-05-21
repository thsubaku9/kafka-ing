#!/bin/zsh

if [[ -z $2 ]]; then
    echo "topic-name cannot be empty"
    exit 1
fi

if [[ "c" == $3 ]]; then
    docker exec -it $1 kafka-topics.sh --create --bootstrap-server kafka:9092 --topic $2
elif [[ "d" == $3]]; then
    docker exec -it $1 kafka-topics.sh --delete --bootstrap-server kafka:9092 --topic $2
fi
