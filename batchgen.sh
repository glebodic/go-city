#!/bin/bash
for i in `seq 1 10`;
do
  echo "Generating" $i "/10"
  go run factory.go
done    