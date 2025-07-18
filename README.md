# Proofpoint URL Decoder

A small Go utility (with an optional HTTP server) to decode URLs rewritten by Proofpoint's URL Defense. Supports v1, v2, and v3 formats.

---

## Features

- **CLI mode**: Decode one or more URLs passed on the command line.
- **Server mode**: Run an HTTP server on port 8089 with a minimal HTML form.
- Supports Proofpoint URL Defense versions **v1**, **v2**, and **v3**.
- Automatically handles HTML and URL escaping.

---

## Requirements

- Go 1.21+
- `github.com/akamensky/argparse`
- (Optional) Any modern web browser for server mode

---

## Installation

```sh
git clone https://github.com/nhdewitt/proofpoint-url-decoder.git
cd proofpoint-decoder
go build -o proofpoint-decoder
```

This will produce the `proofpoint-decoder` binary in the current directory.

---

## Usage

### CLI mode

```sh
./proofpoint-decoder -u "https://urldefense.proofpoint.com/v1/u=https%3A%2F%2Fexample.com%2F&k=ABC123"
```

```
https://example.com/
```

You can pass multiple URLs:

```sh
./proofpoint-decoder -u "<url1>" -u "<url2>" -u "<url3>"
```

### Server Mode

Start the built-in HTTP server on port 8089:

```sh
./proofpoint-decoder -s
```

Open your browser to [http://localhost:8089](http://localhost:8089), paste a URL into the form, and click **Decode**.

---

## Patterns Supported

- **v1**:
  URLs matching `/v1/u=<url-encoded>&k=...`
- **v2**:
  URLs matching `/v2/u=<modified‑base64>&[d|c]=...` (uses `- → %`, `_ → /`)
- **v3**:
  URLs matching `/v3/__<url>__;...!` with embedded Base64‑URL‑encoded token bytes

---

## Acknowledgements

This tool is based on [urldecoder.py](https://help.proofpoint.com/@api/deki/files/2775/urldecoder.py?revision=1) by Eric Van Cleve, licensed under GPLv3.

## License

This project is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).