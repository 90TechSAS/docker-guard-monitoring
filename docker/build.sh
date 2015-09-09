#!/bin/bash

rm -rf dgm
cp -r .. dgm
docker build -t dg-monitoring .
rm -rf dgm
