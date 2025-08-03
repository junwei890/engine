#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose turso $DB_URL up
