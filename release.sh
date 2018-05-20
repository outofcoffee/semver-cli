#!/usr/bin/env bash

set -e
#set -x

CURRENT_VERSION="$( git describe --tags --exact-match )"
GITHUB_USER="outofcoffee"
GITHUB_REPO="semver-cli"

function create_release() {
    CREATE_RELEASE_REQ=$( cat << EOF
{
  "tag_name": "${CURRENT_VERSION}",
  "target_commitish": "master",
  "name": "${CURRENT_VERSION}",
  "body": "Latest edge release",
  "draft": false,
  "prerelease": true
}
EOF
    )

    curl -X POST -u "${GITHUB_USER}:${GITHUB_PASSWORD}" \
        https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases -d "${CREATE_RELEASE_REQ}"
}

function upload_binary() {
    RELEASE_INFO="$( curl -s -u "${GITHUB_USER}:${GITHUB_PASSWORD}" https://api.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases/tags/${CURRENT_VERSION} )"
    RELEASE_ID="$( echo ${RELEASE_INFO} | jq '.id' )"

    curl -X POST -u "${GITHUB_USER}:${GITHUB_PASSWORD}" \
        https://uploads.github.com/repos/${GITHUB_USER}/${GITHUB_REPO}/releases/${RELEASE_ID}/assets?name=$( basename $1 ) \
        -H 'Content-Type: application/octet-stream' \
        --data-binary @$1
}

echo "Building ${CURRENT_VERSION}"

echo "GitHub password for ${GITHUB_USER}:"
read -s GITHUB_PASSWORD

go build .
create_release
upload_binary semver-cli
