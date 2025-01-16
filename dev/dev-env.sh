#!/usr/bin/env sh

echo "Starting!"

echo "Serving Vue frontend first..."
# cd /go/src/github.com/statping-ng/statping-ng
cd frontend && npm run dev &

cd /go/src/github.com/statping-ng/statping-ng modd
cd source
mkdir -p dist
cp -R ../frontend/dist ./dist

echo "Now serving Vue, lets build the golang backend now..."
cd /go/src/github.com/statping-ng/statping-ng
modd -f dev/modd.conf
