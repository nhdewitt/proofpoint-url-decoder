# Proofpoint URL Decoder

A small Go utility (with an optional HTTP server) to decode URLs rewritten by Proofpoint's URL Defense. Supports v1, v2, and v3 formats.

---

## Features

- **CLI mode**: Decode one or more URLs passed on the command line
- **Server mode**: Run an HTTP server with a web interface
- **Mobile-optimized**: Automatic mobile detection with responsive design
- **JSON API**: RESTful endpoint for programmatic access
- **Dark mode**: Toggle between light and dark themes
- **Copy to clipboard**: Click decoded URLs to copy them
- Supports Proofpoint URL Defense versions **v1**, **v2**, and **v3**
- Automatically handles HTML and URL escaping

---

## Requirements

- Go 1.21+
- `github.com/akamensky/argparse`
- (Optional) Any modern web browser for server mode
- (Optional) Docker for containerized deployment

---

## Installation

### Build from Source

```bash
git clone https://github.com/nhdewitt/proofpoint-url-decoder.git
cd proofpoint-url-decoder
go build -o proofpoint-decoder
```

This will produce the `proofpoint-decoder` binary in the current directory.

### Docker Installation

Build the Docker image:

```bash
docker build -t proofpoint-decoder .
```

---

## Usage

### CLI Mode

Decode a single URL:

```bash
./proofpoint-decoder -u "https://urldefense.proofpoint.com/v1/url?u=https%3A%2F%2Fexample.com%2F&k=ABC123"
```

Output:
```
https://example.com/
```

Decode multiple URLs:

```bash
./proofpoint-decoder -u "<url1>" -u "<url2>" -u "<url3>"
```

### Server Mode

#### Native Binary

Start the HTTP server (default port 8089):

```bash
./proofpoint-decoder -s
```

The server will use the port specified in `config.json` if present, or default to 8089.

#### Docker

Run the server in a Docker container:

```bash
# Run once
docker run -p 8089:8089 proofpoint-decoder

# Run in background
docker run -d -p 8089:8089 proofpoint-decoder

# Run with auto-restart on boot
docker run -d --name proofpoint-decoder --restart=unless-stopped -p 8089:8089 proofpoint-decoder
```

#### Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  proofpoint-decoder:
    image: proofpoint-decoder
    container_name: proofpoint-decoder
    restart: unless-stopped
    ports:
      - "8089:8089"
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8089"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

Then run:

```bash
docker-compose up -d
```

### Web Interface

The web interface provides:

- **Desktop view**: Multi-URL form with results display at [http://localhost:8089](http://localhost:8089)
- **Mobile view**: Optimized single-URL interface (auto-detected or at `/m`)
- **Dark mode**: Toggle in the top-right corner
- **Copy functionality**: Click any decoded URL to copy it to clipboard

### JSON API

The server also provides a REST API endpoint:

```bash
curl -X POST http://localhost:8089/api/decode \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "https://urldefense.proofpoint.com/v1/url?u=https%3A%2F%2Fexample.com%2F&k=ABC123",
      "https://urldefense.proofpoint.com/v2/url?u=https-3A__example.org&d=..."
    ]
  }'
```

Response:
```json
{
  "results": ["https://example.com/", "https://example.org"],
  "errors": ["", ""]
}
```

### Docker CLI Mode

You can also use the Docker container in CLI mode:

```bash
docker run --rm proofpoint-decoder ./proofpoint-decoder -u "your-encoded-url-here"
```

---

## Configuration

The server reads configuration from `config.json` if present:

```json
{
  "port": "8089"
}
```

If no config file is found, it defaults to port 8089.

---

## Supported URL Formats

### v1 Format
URLs matching `/v1/u=<url-encoded>&k=...`

### v2 Format  
URLs matching `/v2/u=<modified‑base64>&[d|c]=...`
(uses `- → %`, `_ → /`)

### v3 Format
URLs matching `/v3/__<url>__;...!`
with embedded Base64‑URL‑encoded token bytes

---

## Auto-Start on Boot

To automatically start the Docker container on system boot:

```bash
docker run -d --name proofpoint-decoder --restart=unless-stopped -p 8089:8089 proofpoint-decoder
```

The `--restart=unless-stopped` policy will:
- Restart the container if it crashes
- Restart the container when the system reboots
- Not restart if you manually stop it with `docker stop`

---

## Development

### Project Structure

```
├── main.go              # CLI entry point
├── server.go            # HTTP server setup
├── handlers.go          # HTTP handlers
├── url-defense-decoder.go # Core decoding logic
├── config.json          # Server configuration
├── templates/           # HTML templates
│   ├── form.html        # Desktop form
│   ├── result.html      # Desktop results
│   ├── mobile_form.html # Mobile form
│   └── mobile_result.html # Mobile results
├── static/              # Static assets
│   ├── css/            # Stylesheets
│   └── js/             # JavaScript
└── internal/config/     # Configuration loading
```

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o proofpoint-decoder-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o proofpoint-decoder.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o proofpoint-decoder-macos
```

---

## Acknowledgements

This tool is based on [urldecoder.py](https://help.proofpoint.com/@api/deki/files/2775/urldecoder.py?revision=1) by Eric Van Cleve, licensed under GPLv3.

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).