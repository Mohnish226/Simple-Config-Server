package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var JWT_SECRET = func() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("Warning: JWT_SECRET is not set. Using default 'test' key.")
		secret = "test" // Default value for local testing
	}
	return secret
}()

const BASE_URL = "http://127.0.0.1:8080"
const PRODUCT = "sample"
const ENV = "development"
const CONFIG_KEY = "version"

func base64URLEncode(data []byte) string {
	encoded := base64.URLEncoding.EncodeToString(data)
	return strings.TrimRight(encoded, "=")
}

func generateJWT() string {
	header := `{"alg":"HS256","typ":"JWT"}`
	payload := fmt.Sprintf(`{"exp":%d,"iat":%d}`, time.Now().Add(time.Hour).Unix(), time.Now().Unix())

	headerEncoded := base64URLEncode([]byte(header))
	payloadEncoded := base64URLEncode([]byte(payload))

	data := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, []byte(JWT_SECRET))
	h.Write([]byte(data))
	signature := h.Sum(nil)

	signatureEncoded := base64URLEncode(signature)

	return headerEncoded + "." + payloadEncoded + "." + signatureEncoded
}

func fetchConfig() {
	token := generateJWT()
	url := fmt.Sprintf("%s/%s/%s/%s", BASE_URL, PRODUCT, ENV, CONFIG_KEY)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	if resp.StatusCode == 200 {
		var jsonResponse map[string]interface{}
		json.Unmarshal(body, &jsonResponse)
		fmt.Println("Config Data:", jsonResponse)
	} else {
		fmt.Println("Error:", resp.StatusCode, string(body))
	}
}

func main() {
	if JWT_SECRET == "" {
		fmt.Println("JWT_SECRET environment variable is required")
		fmt.Println("Please set it and try again")
		fmt.Println("Example: export JWT_SECRET=your_secret && go run main.go")
		return
	}
	fetchConfig()
}
