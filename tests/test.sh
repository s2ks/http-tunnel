#!/bin/bash

coproc client { ./client -config client.ini ; }
coproc server { ./server -config server.ini ; }
coproc remote { ./test-remote -listen '127.0.0.123:3211' ; }

trap "kill $client_PID" EXIT
trap "kill $server_PID" EXIT
trap "kill $remote_PID" EXIT

wait $client_PID $server_PID $remote_PID
