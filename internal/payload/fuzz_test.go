package payload

import (
	"strings"
	"testing"
)

// FuzzEncoders checks structural invariants of both encoders under
// arbitrary field values — in particular that escaping prevents
// property/line injection into the vCard envelope. Adopted from the
// 2026-09-05 audit (issue #22).
func FuzzEncoders(f *testing.F) {
	f.Add("ssid", "pass", "WPA", true, "Ada Lovelace")
	f.Add("a;b", `p\q`, "", false, "END:VCARD")
	f.Fuzz(func(t *testing.T, ssid, pass, sec string, hidden bool, name string) {
		if out, err := (WiFi{SSID: ssid, Password: pass, Security: sec, Hidden: hidden}).Encode(); err == nil {
			if !strings.HasPrefix(out, "WIFI:T:") || !strings.HasSuffix(out, ";;") {
				t.Fatalf("malformed wifi payload: %q", out)
			}
		}
		if out, err := (VCard{Name: name, Tel: ssid, Org: pass, Email: sec}).Encode(); err == nil {
			lines := strings.Split(out, "\n")
			if lines[0] != "BEGIN:VCARD" || lines[len(lines)-1] != "END:VCARD" {
				t.Fatalf("malformed vcard envelope: %q", out)
			}
			ends := 0
			for _, l := range lines {
				if l == "END:VCARD" {
					ends++
				}
			}
			if ends != 1 {
				t.Fatalf("line injection: %d END:VCARD lines in %q", ends, out)
			}
		}
	})
}
