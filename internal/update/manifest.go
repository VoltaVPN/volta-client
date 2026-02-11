package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	manifestSchemaVersion = 1
	updatesRootHost       = "updates.voltavpn.com"
	downloadsHost         = "downloads.voltavpn.com"
)

var semverPattern = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

// Manifest describes a signed update announcement for one channel/platform pair.
type Manifest struct {
	ManifestVersion     int    `json:"manifest_version"`
	Channel             string `json:"channel"`
	Platform            string `json:"platform"`
	Arch                string `json:"arch"`
	Version             string `json:"version"`
	ReleaseSeq          uint64 `json:"release_seq"`
	MinSupportedVersion string `json:"min_supported_version"`
	URL                 string `json:"url"`
	SHA256              string `json:"sha256"`
	CreatedAt           string `json:"created_at"`
	ExpiresAt           string `json:"expires_at"`
	KeyID               string `json:"key_id"`
	Signature           string `json:"signature"`
}

// SignedPayload is a canonical subset of Manifest that is covered by signature.
type SignedPayload struct {
	ManifestVersion     int    `json:"manifest_version"`
	Channel             string `json:"channel"`
	Platform            string `json:"platform"`
	Arch                string `json:"arch"`
	Version             string `json:"version"`
	ReleaseSeq          uint64 `json:"release_seq"`
	MinSupportedVersion string `json:"min_supported_version"`
	URL                 string `json:"url"`
	SHA256              string `json:"sha256"`
	CreatedAt           string `json:"created_at"`
	ExpiresAt           string `json:"expires_at"`
	KeyID               string `json:"key_id"`
}

func (m Manifest) ToSignedPayload() SignedPayload {
	return SignedPayload{
		ManifestVersion:     m.ManifestVersion,
		Channel:             m.Channel,
		Platform:            m.Platform,
		Arch:                m.Arch,
		Version:             m.Version,
		ReleaseSeq:          m.ReleaseSeq,
		MinSupportedVersion: m.MinSupportedVersion,
		URL:                 m.URL,
		SHA256:              m.SHA256,
		CreatedAt:           m.CreatedAt,
		ExpiresAt:           m.ExpiresAt,
		KeyID:               m.KeyID,
	}
}

func (m Manifest) ValidateShape(expectedChannel, expectedPlatform, expectedArch string) error {
	if m.ManifestVersion != manifestSchemaVersion {
		return errors.New("unsupported manifest version")
	}

	if m.Channel != expectedChannel || m.Platform != expectedPlatform || m.Arch != expectedArch {
		return errors.New("manifest scope mismatch")
	}

	if !semverPattern.MatchString(m.Version) || !semverPattern.MatchString(m.MinSupportedVersion) {
		return errors.New("invalid semver in manifest")
	}

	if len(m.SHA256) != 64 || !isASCIIHex(m.SHA256) {
		return errors.New("invalid sha256 format")
	}

	if strings.TrimSpace(m.Signature) == "" || strings.TrimSpace(m.KeyID) == "" {
		return errors.New("missing signature fields")
	}

	parsedURL, err := url.Parse(m.URL)
	if err != nil {
		return err
	}
	if parsedURL.Scheme != "https" {
		return errors.New("update url must use https")
	}

	host := strings.ToLower(parsedURL.Hostname())
	if host != downloadsHost && host != updatesRootHost {
		return fmt.Errorf("unexpected update host: %s", host)
	}

	if _, err := time.Parse(time.RFC3339, m.CreatedAt); err != nil {
		return errors.New("invalid created_at")
	}
	if _, err := time.Parse(time.RFC3339, m.ExpiresAt); err != nil {
		return errors.New("invalid expires_at")
	}

	return nil
}

func canonicalPayloadBytes(payload SignedPayload) ([]byte, error) {
	// Struct marshal keeps a deterministic key order for this fixed payload type.
	return json.Marshal(payload)
}

func isASCIIHex(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			continue
		}
		return false
	}
	return true
}
