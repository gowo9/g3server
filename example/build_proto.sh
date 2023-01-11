#!/bin/sh

ProtoDIR=$(dirname $0)/proto
cd $ProtoDIR

buf mod update
buf generate