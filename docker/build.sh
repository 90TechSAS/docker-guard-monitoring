#!/bin/bash

rm -rf dgm
mkdir dgm
cp -r ../config.yaml dgm
cp -r ../core dgm
cp -r ../main.go dgm
cp -r ../transports dgm
cp -r ../utils dgm
docker build -t dg-monitoring .
rm -rf dgm
