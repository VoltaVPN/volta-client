package update

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"
)

func TestVerifyManifest_OK(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	now := time.Now().UTC()
	m := Manifest{
		ManifestVersion:     1,
		Channel:             "stable",
		Platform:            "windows",
		Arch:                "amd64",
		Version:             "1.2.0",
		ReleaseSeq:          12,
		MinSupportedVersion: "1.0.0",
		URL:                 "https://downloads.voltavpn.com/stable/windows/amd64/volta-1.2.0.exe",
		SHA256:              "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		CreatedAt:           now.Add(-time.Minute).Format(time.RFC3339),
		ExpiresAt:           now.Add(time.Hour).Format(time.RFC3339),
		KeyID:               "prod-2026-01",
	}

	payload, err := json.Marshal(m.ToSignedPayload())
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	m.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload))

	err = VerifyManifest(m, map[string]ed25519.PublicKey{"prod-2026-01": pub}, VerifyOptions{
		Channel:  "stable",
		Platform: "windows",
		Arch:     "amd64",
		Now:      now,
		State: State{
			LastSeenReleaseSeq: 11,
			CurrentVersion:     "1.1.9",
		},
	})
	if err != nil {
		t.Fatalf("VerifyManifest error: %v", err)
	}
}

func TestVerifyManifest_DowngradeBlocked(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	now := time.Now().UTC()
	m := Manifest{
		ManifestVersion:     1,
		Channel:             "stable",
		Platform:            "windows",
		Arch:                "amd64",
		Version:             "1.1.0",
		ReleaseSeq:          10,
		MinSupportedVersion: "1.0.0",
		URL:                 "https://downloads.voltavpn.com/stable/windows/amd64/volta-1.1.0.exe",
		SHA256:              "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		CreatedAt:           now.Add(-time.Minute).Format(time.RFC3339),
		ExpiresAt:           now.Add(time.Hour).Format(time.RFC3339),
		KeyID:               "prod-2026-01",
	}
	payload, _ := json.Marshal(m.ToSignedPayload())
	m.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(priv, payload))

	err := VerifyManifest(m, map[string]ed25519.PublicKey{"prod-2026-01": pub}, VerifyOptions{
		Channel:  "stable",
		Platform: "windows",
		Arch:     "amd64",
		Now:      now,
		State: State{
			LastSeenReleaseSeq: 10,
			CurrentVersion:     "1.1.0",
		},
	})
	if err == nil {
		t.Fatal("expected downgrade/replay error")
	}
}

func TestVerifyArtifactSHA256_OK(t *testing.T) {
	data := []byte("volta-artifact")
	sum := sha256.Sum256(data)
	expected := hex.EncodeToString(sum[:])
	if err := VerifyArtifactSHA256(bytes.NewReader(data), expected); err != nil {
		t.Fatalf("VerifyArtifactSHA256 error: %v", err)
	}
}
