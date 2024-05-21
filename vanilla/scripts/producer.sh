#!/bin/zsh

if [[ -z $2 ]]; then
    echo "topic-name cannot be empty"
    exit 1
fi

echo "starting producer"

docker exec -it $1 kafka-console-producer.sh --bootstrap-server kafka:9092 --topic $2

