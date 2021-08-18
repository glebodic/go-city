#!/bin/bash
for i in `seq 1 3`;
do
  echo "Generating" $i "/3"
  go run factory.go
done    