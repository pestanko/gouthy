#! /bin/bash

DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_DATABASE=gouthy
DB_FILES_PATH="db/migrations/psql"
DB_URL="postgres://${DB_USERNAME}:${DB_PASSWORD}@localhost:5432/${DB_DATABASE}?sslmode=disable"

migrate -database $DB_URL -path $DB_FILES_PATH  up


