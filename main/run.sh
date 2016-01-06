#!/bin/bash
echo "Spinning up server..."
go run -race go-server.go &>/dev/null &
sleep 2
echo "Spinning up client..."
go run -race go-client.go
sleep 2
echo "Shutting down server..."
kill `pidof go-server`
echo "Removing database files on disk..."
rm -rf db
exit 0
