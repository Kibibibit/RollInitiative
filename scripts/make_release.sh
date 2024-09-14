#!/bin/bash

V=$1

cd ..

mkdir -p ./build/RollInitiative

go build .

mv roll_initiative ./build/RollInitiative

cp -r srd_data ./build/RollInitiative

cd ./build

zip -r RollInitiative.$V.zip ./RollInitiative

cd ../scripts


