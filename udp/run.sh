#!/bin/sh

# enable jobs
set -m

echo "starting udp server..."
# start server as a background job
go run server/main.go &
echo "udp server started successfully"

# sleep for 2 seconds to ensure server has successfully started
sleep 2

echo "starting udp client..."
go run client/main.go
echo "It's ended before it began"