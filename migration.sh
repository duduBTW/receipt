#!/bin/bash

if [ -z "$1" ]; then
  read -p "Enter the migration name: " migration_name
else
  migration_name="$1"
fi

migrate create -ext sql -dir db/migrations -seq "$migration_name"