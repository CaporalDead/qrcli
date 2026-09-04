// Package payload builds standard QR payload strings from structured
// fields, applying each format's escaping rules so callers (and users)
// don't have to. Like internal/render, everything here is a pure
// function of its inputs and is pinned by exact-string tests.
package payload

import (
	"errors"
	"fmt"
	"strings"
)

// WiFi describes a network in the de-facto standard WIFI: format
// (popularized by ZXing) understood by phone cameras.
type WiFi struct {
	SSID     string
	Password string
	// Security is "WPA", "WEP" or "nopass" (case-insensitive). Empty
	// defaults to WPA when a password is set, nopass otherwise. WPA
	// covers WPA2/WPA3 — scanners treat it as "secured network".
	Security string
	Hidden   bool
}

// wifiEscaper escapes the characters reserved by the WIFI: format.
// A Replacer substitutes all patterns in a single pass, so inserted
// backslashes are never re-escaped.
var wifiEscaper = strings.NewReplacer(
	`\`, `\\`, `;`, `\;`, `,`, `\,`, `:`, `\:`, `"`, `\"`,
)

// Encode returns the WIFI:T:...;S:...;; payload.
func (w WiFi) Encode() (string, error) {
	if w.SSID == "" {
		return "", errors.New("wifi: --ssid is required")
	}
	security := w.Security
	if security == "" {
		if w.Password == "" {
			security = "nopass"
		} else {
			security = "WPA"
		}
	}
	switch strings.ToLower(security) {
	case "wpa":
		security = "WPA"
	case "wep":
		security = "WEP"
	case "nopass":
		security = "nopass"
	default:
		return "", fmt.Errorf("wifi: invalid security type %q (want WPA, WEP or nopass)", w.Security)
	}
	if security == "nopass" && w.Password != "" {
		return "", errors.New("wifi: --pass conflicts with --type nopass")
	}
	if security != "nopass" && w.Password == "" {
		return "", fmt.Errorf("wifi: --pass is required with --type %s", security)
	}

	var b strings.Builder
	b.WriteString("WIFI:T:")
	b.WriteString(security)
	b.WriteString(";S:")
	b.WriteString(wifiEscaper.Replace(w.SSID))
	b.WriteString(";")
	if w.Password != "" {
		b.WriteString("P:")
		b.WriteString(wifiEscaper.Replace(w.Password))
		b.WriteString(";")
	}
	if w.Hidden {
		b.WriteString("H:true;")
	}
	b.WriteString(";")
	return b.String(), nil
}

// VCard describes a contact card encoded as vCard 3.0 — chosen over the
// more compact MECARD for its universal scanner support (issue #17).
type VCard struct {
	// Name is the full display name (FN). The structured N property is
	// derived with a best-effort "last word is the family name" split;
	// scanners display FN, so the heuristic only affects how fields map
	// when the contact is saved.
	Name  string
	Tel   string
	Email string
	Org   string
	URL   string
}

// vcardEscaper escapes the characters reserved by vCard 3.0 text values
// (RFC 2426): backslash, semicolon, comma, and raw newlines.
var vcardEscaper = strings.NewReplacer(
	`\`, `\\`, `;`, `\;`, `,`, `\,`, "\r\n", `\n`, "\n", `\n`,
)

// Encode returns the BEGIN:VCARD...END:VCARD payload. Lines are joined
// with LF: the spec mandates CRLF, but scanners accept LF and it keeps
// the QR code smaller (issue #17).
func (v VCard) Encode() (string, error) {
	name := strings.TrimSpace(v.Name)
	if name == "" {
		return "", errors.New("vcard: --name is required")
	}
	family, given := splitName(name)
	lines := []string{
		"BEGIN:VCARD",
		"VERSION:3.0",
		"N:" + vcardEscaper.Replace(family) + ";" + vcardEscaper.Replace(given),
		"FN:" + vcardEscaper.Replace(name),
	}
	if v.Org != "" {
		lines = append(lines, "ORG:"+vcardEscaper.Replace(v.Org))
	}
	if v.Tel != "" {
		lines = append(lines, "TEL:"+vcardEscaper.Replace(v.Tel))
	}
	if v.Email != "" {
		lines = append(lines, "EMAIL:"+vcardEscaper.Replace(v.Email))
	}
	if v.URL != "" {
		lines = append(lines, "URL:"+vcardEscaper.Replace(v.URL))
	}
	lines = append(lines, "END:VCARD")
	return strings.Join(lines, "\n"), nil
}

// splitName maps a display name onto vCard's structured N property.
func splitName(full string) (family, given string) {
	fields := strings.Fields(full)
	if len(fields) == 1 {
		return fields[0], ""
	}
	return fields[len(fields)-1], strings.Join(fields[:len(fields)-1], " ")
}
