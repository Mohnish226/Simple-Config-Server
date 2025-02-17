import base64
import hashlib
import hmac
import http.client
import json
import os
import time

JWT_SECRET = os.getenv("JWT_SECRET")
BASE_URL = "127.0.0.1:8080"
PRODUCT = "sample"
ENV = "development"
CONFIG_KEY = "version"


def base64_url_encode(data):
    """Encodes data in Base64 URL-safe format without padding."""
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def generate_jwt():
    """Manually creates a JWT token (HS256) without using `jwt` package."""
    header = {"alg": "HS256", "typ": "JWT"}
    payload = {"exp": int(time.time()) + 3600, "iat": int(time.time())}

    header_encoded = base64_url_encode(json.dumps(header).encode())
    payload_encoded = base64_url_encode(json.dumps(payload).encode())

    signature = hmac.new(
        JWT_SECRET.encode(),
        f"{header_encoded}.{payload_encoded}".encode(),
        hashlib.sha256,
    ).digest()

    signature_encoded = base64_url_encode(signature)
    return f"{header_encoded}.{payload_encoded}.{signature_encoded}"


def fetch_config():
    """Sends an HTTP GET request to fetch config data."""
    token = generate_jwt()
    url = f"/{PRODUCT}/{ENV}/{CONFIG_KEY}"

    conn = http.client.HTTPConnection(BASE_URL)

    headers = {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

    conn.request("GET", url, headers=headers)
    response = conn.getresponse()

    response_data = response.read().decode()
    conn.close()

    if response.status == 200:
        print("Config Data:", json.loads(response_data))
    else:
        print("Error:", response.status, response_data)


if __name__ == "__main__":
    if JWT_SECRET is None:
        print("JWT_SECRET environment variable is not set.")
        print("Please set the JWT_SECRET environment variable which is used to sign JWT tokens.")
        print("Example Usage: export JWT_SECRET=your_secret_key && python3 python-client.py")
        exit(1)
    fetch_config()
