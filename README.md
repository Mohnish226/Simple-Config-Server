# Simple-Config-Server

A lightweight configuration management service that loads YAML-based configurations from a structured directory and exposes them via an HTTP API. It also includes authentication, IP filtering, and rate limiting.

### Use Cases

- Centralized configuration management for microservices.
- Securely store and retrieve configurations for different environments (e.g., development, staging, production).
- Fetch configurations from a remote server using a simple API.

### Illustration

![Simple-Config-Server](extra/usecase.drawio.png)

### Project Structure

```
Simple-Config-Server
 │
 │── /configurations            # Stores project-specific configuration files
 │   ├── /sample
 │   │    └── development.yml   # Example configuration file
 │   └── Readme.md              # Documentation for adding configurations
 │
 │── /internal                  # Internal modules for core functionality
 │   ├── /auth                  # JWT-based authentication
 │   │    └── jwt.go
 │   │
 │   ├── /config                # Configuration loader & file watcher
 │   │    ├── config.go
 │   │    └── watcher.go
 │   │
 │   ├── /handler               # API handlers for retrieving configurations
 │   │    └── handler.go
 │   │
 │   ├── /ipfilter              # IP whitelisting for security
 │   │    ├── filter.go
 │   │    └── watcher.go
 │   │
 │   ├── /logger                # Logging utility
 │   │    └── logger.go
 │   │
 │   ├── /rate_limiter          # Rate limiting middleware
 │   │    └── limiter.go
 │   │
 │   └── /scaffolding           # Create the Configurations directory structure
 │        └── scaffold.go
 │
 │── /clients                   # Example clients to fetch configurations
 │   ├── golang-client.go       # Example client in Go
 |   └─ python-client.py        # Example client in Python
 │
 │── .gitignore                 # Git ignored files
 │── allowed_ips.txt            # List of allowed IPs for access control
 │── allowed_ips.txt.example    # Example IP allowlist
 │── application.log            # Log file
 │── go.mod                     # Go module dependencies
 │── go.sum                     # Go module checksum file
 │── main.go                    # Entry point of the application
 │── LICENSE                    # License file
 └── README.md                  # Documentation
```

### Configuration Files

Please refer to the [configurations](configurations/Readme.md) documentation for adding configuration files.

### IP Allowlist Configuration

The `allowed_ips.txt` file supports comments and can be organized with sections:

```
# Local development
127.0.0.1
::1

# Production servers
10.0.0.0/8
172.16.0.0/12

# Staging environment
192.168.1.0/24
```

- Lines starting with `#` are treated as comments
- Inline comments (after `#`) are supported
- Empty lines are ignored
- If the file is empty, all IPs are allowed

### Usage

1. Clone the repository
2. Build the application:
    ```bash
    make build
    ```
3. Run the application:
    ```bash
    make run
    ```

The application will automatically:
- Create the configurations directory if it doesn't exist
- Create a sample configuration file in `configurations/sample/development.yml`
- Create an empty `allowed_ips.txt` file if it doesn't exist
- All files will be created in the same directory as the binary

Or with custom configuration paths:
```bash
./bin/simple-config-server --config-dir=/path/to/configs --allowed-ips=/path/to/ips.txt
```

Environment variables can also be used:
```bash
export CONFIG_DIR=/path/to/configs
export ALLOWED_IPS_FILE=/path/to/ips.txt
export PORT=8080
export JWT_SECRET=secret
./bin/simple-config-server
```

4. Access the API:
    ```bash
    curl -H "Authorization: Bearer <your_token>" -X GET http://127.0.0.1:8080/<project>/<environment>/<config>
    ```

### Build Client to Fetch Configurations

Please refer to the example client code in the [client](clients) directory.

### Planned Features 🚀

- [ ] Support additional configuration formats (e.g., JSON, TOML) for greater flexibility.
- [ ] Enable configuration push to allow updates directly from clients.
- [ ] Introduce versioning to track and manage configuration changes.
- [ ] Implement encryption & decryption to enhance configuration security.

> Note: The configuration file should not contain any sensitive information such as passwords, API keys, etc. Sensitive information should be stored in a secure location and accessed using environment variables. This project is intended for use with non-sensitive configuration settings only.