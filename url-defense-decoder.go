package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

type urlDefenseDecoder struct {
	udPattern      *regexp.Regexp
	v1Pattern      *regexp.Regexp
	v2Pattern      *regexp.Regexp
	v3Pattern      *regexp.Regexp
	v3TokenPattern *regexp.Regexp
	v3SingleSlash  *regexp.Regexp
	v3RunMapping   map[rune]int
}

func buildV3RunMapping() map[rune]int {
	runValues := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	runLength := 2
	v3RunMap := make(map[rune]int)

	for _, r := range runValues {
		v3RunMap[r] = runLength
		runLength++
	}

	return v3RunMap
}

func (d *urlDefenseDecoder) Decode(rewrittenURL string) (string, error) {
	if parts := d.udPattern.FindStringSubmatch(rewrittenURL); parts != nil {
		switch parts[1] {
		case "v1":
			return d.decodeV1(rewrittenURL)
		case "v2":
			return d.decodeV2(rewrittenURL)
		case "v3":
			return d.decodeV3(rewrittenURL)
		}
	}
	return "", errors.New("does not appear to be a URL Defense URL")
}

func (d *urlDefenseDecoder) decodeV1(rewrittenURL string) (string, error) {
	parts := d.v1Pattern.FindStringSubmatch(rewrittenURL)
	if parts == nil {
		return "", errors.New("error parsing URL")
	}
	urlEncoded := parts[1]

	htmlEncoded, err := url.QueryUnescape(urlEncoded)
	if err != nil {
		return "", err
	}

	decoded := html.UnescapeString(htmlEncoded)
	return decoded, nil
}

func (d *urlDefenseDecoder) decodeV2(rewrittenURL string) (string, error) {
	parts := d.v2Pattern.FindStringSubmatch(rewrittenURL)
	if parts == nil {
		return "", errors.New("error parsing URL")
	}

	var v2Replacer = strings.NewReplacer(
		"-", "%",
		"_", "/",
	)

	urlSpecialEncoded := v2Replacer.Replace(parts[1])

	htmlEncoded, err := url.QueryUnescape(urlSpecialEncoded)
	if err != nil {
		return "", err
	}

	decoded := html.UnescapeString(htmlEncoded)
	return decoded, nil
}

func (d *urlDefenseDecoder) decodeV3(rewrittenURL string) (string, error) {
	parts := d.v3Pattern.FindStringSubmatch(rewrittenURL)
	if parts == nil {
		return "", fmt.Errorf("decodev3: no v3 match in %q", rewrittenURL)
	}
	names := d.v3Pattern.SubexpNames()

	var urlMatch, encBytes string
	for i, name := range names {
		switch name {
		case "url":
			urlMatch = parts[i]
		case "enc_bytes":
			encBytes = parts[i] + "=="
		}
	}

	if ss := d.v3SingleSlash.FindStringSubmatch(urlMatch); len(ss) >= 3 {
		urlMatch = ss[1] + "/" + ss[2]
	}

	encodedUrl, err := url.QueryUnescape(urlMatch)
	if err != nil {
		return "", fmt.Errorf("decodev3: unquote: %w", err)
	}

	decBytes, err := base64.URLEncoding.DecodeString(encBytes)
	if err != nil {
		return "", fmt.Errorf("decodev3: base64 decode: %w", err)
	}
	currentMarker := 0

	replaceToken := func(token string) string {
		if token == "*" {
			b := decBytes[currentMarker]
			currentMarker++
			return string(b)
		}
		if strings.HasPrefix(token, "**") {
			last := rune(token[len(token)-1])
			runLen, ok := d.v3RunMapping[last]
			if !ok {
				return ""
			}
			end := min(currentMarker+runLen, len(decBytes))
			chunk := decBytes[currentMarker:end]
			currentMarker += runLen
			return string(chunk)
		}
		return token
	}

	var substituteTokens func(text string, pos int) string
	substituteTokens = func(text string, pos int) string {
		suffix := text[pos:]
		loc := d.v3TokenPattern.FindStringIndex(suffix)
		if loc == nil {
			return text[pos:]
		}
		start, end := pos+loc[0], pos+loc[1]
		prefix := text[pos:start]
		token := text[start:end]
		return prefix + replaceToken(token) + substituteTokens(text, end)
	}

	return substituteTokens(encodedUrl, 0), nil
}
