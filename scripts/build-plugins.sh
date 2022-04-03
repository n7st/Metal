#!/usr/bin/env bash

shopt -s extglob;

for file in ./internal/app/plugins/**/!(*_test).go;
do
    PLUGIN=$(basename $(dirname $file));
    go build -buildmode=plugin -o ./plugins/${PLUGIN}.so $file;
done
