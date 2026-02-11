package update

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strings"
	"time"
)

const (
	maxClockSkew          = 5 * time.Minute
	maxArtifactBytes int64 = 1 << 30 // 1 GiB hard guard
)

type State struct {
	LastSeenReleaseSeq uint64
	CurrentVersion     string
}

type VerifyOptions struct {
	Channel  string
	Platform string
	Arch     string
	Now      time.Time
	State    State
}

func VerifyManifest(m Manifest, keyring map[string]ed25519.PublicKey, opts VerifyOptions) error {
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}

	if err := m.ValidateShape(opts.Channel, opts.Platform, opts.Arch); err != nil {
		return err
	}

	pub, ok := keyring[m.KeyID]
	if !ok {
		return errors.New("unknown key id")
	}

	createdAt, _ := time.Parse(time.RFC3339, m.CreatedAt)
	expiresAt, _ := time.Parse(time.RFC3339, m.ExpiresAt)
	if opts.Now.Before(createdAt.Add(-maxClockSkew)) || opts.Now.After(expiresAt.Add(maxClockSkew)) {
		return errors.New("manifest is outside valid time window")
	}

	if m.ReleaseSeq <= opts.State.LastSeenReleaseSeq {
		return errors.New("downgrade/replay by release sequence")
	}

	if strings.TrimSpace(opts.State.CurrentVersion) != "" {
		compare, err := compareSemver(m.Version, opts.State.CurrentVersion)
		if err != nil {
			return err
		}
		if compare <= 0 {
			return errors.New("update version must be greater than current version")
		}
	}

	payloadBytes, err := canonicalPayloadBytes(m.ToSignedPayload())
	if err != nil {
		return err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(m.Signature)
	if err != nil {
		return errors.New("invalid signature encoding")
	}
	if !ed25519.Verify(pub, payloadBytes, signatureBytes) {
		return errors.New("invalid manifest signature")
	}

	return nil
}

func VerifyArtifactSHA256(r io.Reader, expectedHex string) error {
	expectedHex = strings.ToLower(strings.TrimSpace(expectedHex))
	if len(expectedHex) != 64 || !isASCIIHex(expectedHex) {
		return errors.New("invalid expected sha256")
	}

	lr := &io.LimitedReader{R: r, N: maxArtifactBytes + 1}
	h := sha256.New()
	if _, err := io.Copy(h, lr); err != nil {
		return err
	}
	if lr.N <= 0 {
		return errors.New("artifact exceeds max size")
	}

	sum := h.Sum(nil)
	if hex.EncodeToString(sum) != expectedHex {
		return errors.New("artifact hash mismatch")
	}

	return nil
}

func compareSemver(a, b string) (int, error) {
	ap, err := parseSemver(a)
	if err != nil {
		return 0, err
	}
	bp, err := parseSemver(b)
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ {
		if ap[i] > bp[i] {
			return 1, nil
		}
		if ap[i] < bp[i] {
			return -1, nil
		}
	}
	return 0, nil
}

func parseSemver(s string) ([3]int, error) {
	var out [3]int
	if !semverPattern.MatchString(s) {
		return out, errors.New("invalid semver")
	}

	parts := strings.Split(s, ".")
	for i := range parts {
		n := 0
		for j := 0; j < len(parts[i]); j++ {
			n = n*10 + int(parts[i][j]-'0')
		}
		out[i] = n
	}
	return out, nil
}
