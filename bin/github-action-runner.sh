#!/bin/bash

source ./.env

REG_TOKEN=$(curl -X POST -H "Authorization: token ${GITHUP_ACCESS_TOKEN}" -H "Accept: application/vnd.github+json" https://api.github.com/repos/${GITHUB_REPO}/actions/runners/registration-token | jq .token --raw-output)

cd ./actions-runner

./config.sh --url https://github.com/${GITHUB_REPO} --token ${REG_TOKEN}

cleanup() {
    echo "Removing runner..."
    ./config.sh remove --unattended --token ${REG_TOKEN}
}

trap 'cleanup; exit 130' INT
trap 'cleanup; exit 143' TERM

./run.sh & wait $!
