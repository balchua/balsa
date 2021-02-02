#!/bin/sh

test_loop() {
	for i in $(seq 1000); do
		OUT=$(curl -X GET "http://localhost:8000/set?key=${i}&value=${i}") || sleep 2s
	done
	echo "Loop ${1} exited"
}

test_loop