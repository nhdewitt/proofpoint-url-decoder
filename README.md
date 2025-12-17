# Proofpoint URL Decoder

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Docker Support](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)
![License](https://img.shields.io/badge/License-GPLv3-blue.svg)

A robust Go utility and HTTP server designed to decode URLs rewritten by Proofpoint's URL Defense system. It supports **v1**, **v2**, and **v3** formats and includes a modern, responsive web interface.

---

## 📸 Screenshots

| **Desktop View** | **Mobile View** |
|:---:|:---:|
| *[Place Desktop Screenshot Here]* | *[Place Mobile Screenshot Here]* |
| *Clean, multi-line decoding with Dark Mode* | *Touch-optimized interface* |

---

## ✨ Features

- **🚀 CLI Mode**: Decode URLs directly from your terminal.
- **🌐 Server Mode**: Fast HTTP server with a web interface.
- **📱 Mobile-First**: Automatically detects mobile devices and serves a touch-optimized UI.
- **🌙 Dark Mode**: Built-in toggle for light and dark themes.
- **📋 One-Click Copy**: Click any decoded result to instantly copy it to your clipboard.
- **🤖 JSON API**: RESTful endpoint for programmatic integration.
- **🛡️ Full Support**: Handles Proofpoint URL Defense **v1**, **v2**, and **v3** formats seamlessly.

---

## 🛠️ Requirements

- **Go 1.21+** (for building from source)
- **Docker** (optional, for containerized deployment)

---

## 📦 Installation

### Option 1: Build from Source

```bash
git clone [https://github.com/nhdewitt/proofpoint-url-decoder.git](https://github.com/nhdewitt/proofpoint-url-decoder.git)
cd proofpoint-url-decoder
go build -o proofpoint-decoder
```
*This generates a `proofpoint-decoder` binary in your current directory.*

### Option 2: Docker

Build the image locally:

```bash
docker build -t proofpoint-decoder .
```

---

## 🚀 Usage

### 1. CLI Mode

**Decode a single URL:**
```bash
./proofpoint-decoder -u "[https://urldefense.proofpoint.com/v1/url?u=https%3A%2F%2Fexample.com&k=ABC](https://urldefense.proofpoint.com/v1/url?u=https%3A%2F%2Fexample.com&k=ABC)..."
# Output: [https://example.com](https://example.com)
```

**Decode multiple URLs:**
```bash
./proofpoint-decoder -u "<url1>" -u "<url2>" -u "<url3>"
```

### 2. Server Mode

Start the server (defaults to port `8089`):
```bash
./proofpoint-decoder -s
```
*Access the web interface at [http://localhost:8089](http://localhost:8089)*

---

## 🐧 Running as a Linux Service (Systemd)

If you prefer not to use Docker, you can run the decoder as a background service using `systemd`.

1. **Install binary and assets**:
   Since the application requires `templates/` and `static/` directories to run, we must place them in a working directory.

   ```bash
   # 1. Move binary to path
   sudo mv proofpoint-decoder /usr/local/bin/
   sudo chmod +x /usr/local/bin/proofpoint-decoder

   # 2. Create working directory and copy assets
   sudo mkdir -p /var/lib/proofpoint-url-decoder
   sudo cp -r templates static config.json /var/lib/proofpoint-url-decoder/
   ```

2. **Create the service file**:
   ```bash
   sudo nano /etc/systemd/system/proofpoint-decoder.service
   ```

3. **Paste the following configuration**:
   ```ini
   [Unit]
   Description=Proofpoint URL Decoder Server
   After=network.target

   [Service]
   Type=simple
   User=root
   # Point to where we copied the templates/static folders
   WorkingDirectory=/var/lib/proofpoint-url-decoder
   ExecStart=/usr/local/bin/proofpoint-decoder -s
   Restart=always
   RestartSec=5s

   [Install]
   WantedBy=multi-user.target
   ```

4. **Enable and Start the service**:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable --now proofpoint-decoder
   ```

5. **Check status**:
   ```bash
   sudo systemctl status proofpoint-decoder
   ```

---

## 🐳 Docker Usage

**Run in background (Standard):**
```bash
docker run -d -p 8089:8089 --name pp-decoder proofpoint-decoder
```

**Run with auto-restart (Recommended for Servers):**
```bash
docker run -d \
  --name pp-decoder \
  --restart=unless-stopped \
  -p 8089:8089 \
  proofpoint-decoder
```

**Run as a one-off CLI tool:**
```bash
docker run --rm proofpoint-decoder ./proofpoint-decoder -u "your-encoded-url"
```

### Docker Compose

Create a `docker-compose.yml`:

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
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8089"]
      interval: 30s
      timeout: 10s
      retries: 3
```

---

## 🔌 API Reference

The server exposes a JSON endpoint at `POST /api/decode`.

**Request:**
```bash
curl -X POST http://localhost:8089/api/decode \
  -H "Content-Type: application/json" \
  -d '{
    "urls": [
      "[https://urldefense.proofpoint.com/v2/url?u=https-3A__example.com](https://urldefense.proofpoint.com/v2/url?u=https-3A__example.com)...",
      "[https://invalid-url.com](https://invalid-url.com)"
    ]
  }'
```

**Response:**
```json
{
  "results": [
    "[https://example.com](https://example.com)",
    ""
  ],
  "errors": [
    "",
    "error: invalid proofpoint format"
  ]
}
```

---

## ⚙️ Configuration

The application looks for a `config.json` file in the working directory. If not found, it defaults to port `8089`.

```json
{
  "port": "8089"
}
```

---

## 🏗️ Development

### Project Structure

```text
.
├── main.go                # CLI entry point
├── server.go              # HTTP server setup
├── handlers.go            # HTTP handlers
├── url-defense-decoder.go # Core decoding logic
├── config.json            # Configuration
├── Dockerfile             # Container definition
├── templates/             # HTML Templates
│   ├── form.html          # Desktop UI
│   ├── mobile_form.html   # Mobile UI
│   └── ...
└── static/                # Assets (CSS/JS)
    ├── css/
    └── js/
```

### Build & Test

```bash
# Run Unit Tests
go test ./...

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o proofpoint-decoder-linux

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o proofpoint-decoder.exe
```

---

## 📜 Acknowledgements & License

**Original Logic:** Based on [urldecoder.py](https://help.proofpoint.com/@api/deki/files/2775/urldecoder.py?revision=1) by Eric Van Cleve.

**License:** This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).