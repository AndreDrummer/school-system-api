.PHONY: dafault build test run docs clean

APP_NAME=school-system-api

default: run

run:
	@go run cmd/main.go