#!/bin/bash
go build -o tiny main/tiny_tool.go

tar -cvf tiny.tar.gz tiny

shasum -a 256 tiny.tar.gz

