#!/usr/bin/env bash

for file in ./internal/app/plugins/**/*.go;
do
    PLUGIN=$(basename $(dirname $file));
    go build -buildmode=plugin -o ./plugins/${PLUGIN}.so $file;
done
