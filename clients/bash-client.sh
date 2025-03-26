#!/bin/bash

# Default JWT_SECRET if not set
JWT_SECRET="${JWT_SECRET:-test}"

# Constants
BASE_URL="http://127.0.0.1:8080"
PRODUCT="sample"
ENV="development"
CONFIG_KEY="version"

# Base64 URL encoding function
base64url_encode() {
    echo -n "$1" | base64 | tr '+/' '-_' | tr -d '='
}

# Generate JWT Token
generate_jwt() {
    header='{"alg":"HS256","typ":"JWT"}'
    payload="{\"exp\":$(($(date +%s) + 3600)),\"iat\":$(date +%s)}"

    header_encoded=$(base64url_encode "$header")
    payload_encoded=$(base64url_encode "$payload")

    data="$header_encoded.$payload_encoded"
    signature=$(echo -n "$data" | openssl dgst -sha256 -hmac "$JWT_SECRET" -binary | base64 | tr '+/' '-_' | tr -d '=')

    echo "$data.$signature"
}

# Fetch config using the JWT token
fetch_config() {
    token=$(generate_jwt)
    url="$BASE_URL/$PRODUCT/$ENV/$CONFIG_KEY"

    response=$(curl -s -w "%{http_code}" -X GET "$url" -H "Authorization: Bearer $token" -H "Content-Type: application/json")

    response_code=$(echo "$response" | tail -c 4)

    if [[ "$response_code" == "200" ]]; then
        echo "Config Data:"
        echo "${response:0:${#response}-3}"
    else
        echo "Error: HTTP $response"
    fi
}

# Main script execution
if [ -z "$JWT_SECRET" ]; then
    echo "JWT_SECRET environment variable is required."
    echo "Please set it and try again."
    echo "Example: export JWT_SECRET=your_secret && ./bash-client.sh"
    exit 1
fi

fetch_config