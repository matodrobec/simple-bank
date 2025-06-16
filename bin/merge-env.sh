#!/bin/bash

while IFS='=' read -r key value; do
  if grep -q "^$key=" app.env; then
    sed -i "s|^$key=.*|$key=$value|" app.env
  else
    echo "$key=$value" >> app.env
  fi
done < app.env.prod
