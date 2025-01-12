#!/bin/bash

# Wait for the API to be available
echo "Waiting for API to be ready..."
until curl -s http://api:8080/task >/dev/null; do
  sleep 1
done

# Add default tasks
echo "Adding default tasks..."
curl --location 'localhost:8080/task' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Task 1"
}'

curl --location 'localhost:8080/task' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Task 2"
}'

echo "Default tasks added."
