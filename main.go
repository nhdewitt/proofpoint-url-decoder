package main

import (
	"fmt"
	"html"
	"os"
	"regexp"

	"github.com/akamensky/argparse"
	"github.com/nhdewitt/proofpoint-url-decoder/internal/config"
)

func main() {
	udd := urlDefenseDecoder{
		udPattern:      regexp.MustCompile(`^https://urldefense(?:\.proofpoint)?\.com/(v[0-9])/$`),
		v1Pattern:      regexp.MustCompile(`u=(?P<url>.+?)&k=`),
		v2Pattern:      regexp.MustCompile(`u=(?P<url>.+?)&[dc]=`),
		v3Pattern:      regexp.MustCompile(`v3/__(?P<url>.+?)__;(?P<enc_bytes>.*?)!`),
		v3TokenPattern: regexp.MustCompile(`\*(\*.)?`),
		v3SingleSlash:  regexp.MustCompile(`^(?i)([a-z0-9+.-]+:/)([^/].+)`),
		v3RunMapping:   buildV3RunMapping(),
	}

	parser := argparse.NewParser("proofpoint-decoder", "Decode URLs rewritten by URL Defense. Supports v1, v2, and v3 URLs.")

	serverMode := parser.Flag("s", "server",
		&argparse.Options{
			Required: false,
			Help:     "Run in HTTP-server mode (no URLs on CLI)",
		})

	urls := parser.StringList("u", "urls",
		&argparse.Options{
			Required: false,
			Help:     "one or more rewritten URLs to decode",
		},
	)

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *serverMode {
		c, err := config.LoadConfig("config.json")
		if err != nil {
			fmt.Errorf("error opening config file: %w", err)
		}
		runServer(&udd, c)
		return
	}

	if len(*urls) == 0 {
		fmt.Println("Error: must either use -s for serer mode or include at least one url with -u/--urls")
		fmt.Print(parser.Usage(nil))
		os.Exit(1)
	}

	for _, raw := range *urls {
		u := html.UnescapeString(raw)
		decoded, err := udd.Decode(u)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to decode %q: %v\n", u, err)
			continue
		}
		fmt.Println(decoded)
	}
}
