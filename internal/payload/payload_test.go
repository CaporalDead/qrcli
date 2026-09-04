package payload

import (
	"strings"
	"testing"
)

func TestWiFiEncode(t *testing.T) {
	tests := []struct {
		name string
		in   WiFi
		want string
	}{
		{
			name: "password defaults to WPA",
			in:   WiFi{SSID: "Home", Password: "hunter2"},
			want: "WIFI:T:WPA;S:Home;P:hunter2;;",
		},
		{
			name: "no password defaults to nopass, P omitted",
			in:   WiFi{SSID: "FreeSpot"},
			want: "WIFI:T:nopass;S:FreeSpot;;",
		},
		{
			name: "hidden network",
			in:   WiFi{SSID: "Home", Password: "x", Hidden: true},
			want: "WIFI:T:WPA;S:Home;P:x;H:true;;",
		},
		{
			name: "explicit WEP, case-insensitive",
			in:   WiFi{SSID: "Old", Password: "k", Security: "wep"},
			want: "WIFI:T:WEP;S:Old;P:k;;",
		},
		{
			name: "explicit nopass, case-insensitive",
			in:   WiFi{SSID: "Open", Security: "NOPASS"},
			want: "WIFI:T:nopass;S:Open;;",
		},
		{
			name: "reserved characters are escaped",
			in:   WiFi{SSID: `a;b:c"d,e\f`, Password: "p;q"},
			want: `WIFI:T:WPA;S:a\;b\:c\"d\,e\\f;P:p\;q;;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Encode()
			if err != nil {
				t.Fatalf("Encode() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWiFiEncodeErrors(t *testing.T) {
	tests := []struct {
		name   string
		in     WiFi
		wantIn string
	}{
		{"missing ssid", WiFi{Password: "x"}, "--ssid is required"},
		{"invalid security", WiFi{SSID: "x", Password: "p", Security: "wpa9"}, "invalid security"},
		{"pass with nopass", WiFi{SSID: "x", Password: "p", Security: "nopass"}, "conflicts"},
		{"WPA without pass", WiFi{SSID: "x", Security: "WPA"}, "--pass is required"},
		{"WEP without pass", WiFi{SSID: "x", Security: "wep"}, "--pass is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.in.Encode()
			if err == nil {
				t.Fatal("Encode() expected an error")
			}
			if !strings.Contains(err.Error(), tt.wantIn) {
				t.Errorf("error %q does not contain %q", err, tt.wantIn)
			}
		})
	}
}

func TestVCardEncode(t *testing.T) {
	tests := []struct {
		name string
		in   VCard
		want string
	}{
		{
			name: "all fields",
			in: VCard{
				Name:  "Ada Lovelace",
				Tel:   "+44 20 7946 0958",
				Email: "ada@example.org",
				Org:   "Analytical Engines Ltd",
				URL:   "https://example.org",
			},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:Lovelace;Ada\nFN:Ada Lovelace\n" +
				"ORG:Analytical Engines Ltd\nTEL:+44 20 7946 0958\n" +
				"EMAIL:ada@example.org\nURL:https://example.org\nEND:VCARD",
		},
		{
			name: "single-word name goes to the family slot",
			in:   VCard{Name: "Ada"},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:Ada;\nFN:Ada\nEND:VCARD",
		},
		{
			name: "last word is the family name",
			in:   VCard{Name: "Jean Luc Picard"},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:Picard;Jean Luc\nFN:Jean Luc Picard\nEND:VCARD",
		},
		{
			name: "surrounding whitespace is trimmed",
			in:   VCard{Name: "  Ada  "},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:Ada;\nFN:Ada\nEND:VCARD",
		},
		{
			name: "reserved characters are escaped, colon is not",
			in:   VCard{Name: "A; B, C", URL: "https://example.org/a,b"},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:C;A\\; B\\,\nFN:A\\; B\\, C\n" +
				"URL:https://example.org/a\\,b\nEND:VCARD",
		},
		{
			name: "raw newline in a field becomes a literal backslash-n",
			in:   VCard{Name: "Ada", Org: "Line1\nLine2"},
			want: "BEGIN:VCARD\nVERSION:3.0\nN:Ada;\nFN:Ada\nORG:Line1\\nLine2\nEND:VCARD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.Encode()
			if err != nil {
				t.Fatalf("Encode() unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Encode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestVCardEncodeErrors(t *testing.T) {
	for _, in := range []VCard{{}, {Name: "   "}, {Name: "\n"}} {
		if _, err := in.Encode(); err == nil {
			t.Errorf("Encode() with name %q expected an error", in.Name)
		}
	}
}

func TestSplitName(t *testing.T) {
	tests := []struct {
		full, family, given string
	}{
		{"Ada", "Ada", ""},
		{"Ada Lovelace", "Lovelace", "Ada"},
		{"Jean Luc Picard", "Picard", "Jean Luc"},
	}
	for _, tt := range tests {
		family, given := splitName(tt.full)
		if family != tt.family || given != tt.given {
			t.Errorf("splitName(%q) = (%q, %q), want (%q, %q)",
				tt.full, family, given, tt.family, tt.given)
		}
	}
}
