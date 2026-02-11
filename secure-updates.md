# VoltaVPN Secure Updates (design)

This document defines release signing and client-side verification for desktop updates.
Primary domain: `voltavpn.com`.

## Update channels

- Manifest endpoint: `https://updates.voltavpn.com/<channel>/<platform>/<arch>/manifest.json`
- Artifact endpoint: `https://downloads.voltavpn.com/<channel>/<platform>/<arch>/<artifact>`
- Transport: HTTPS only, no redirects to a different host.

## Signature model

- Signature algorithm: Ed25519.
- Client contains a pinned public keyring (`key_id -> public key`).
- Manifest is signed as detached signature over canonical payload (all fields except `signature`).
- Key rotation: ship at least two keys during transition windows.

## Manifest format (v1)

```json
{
  "manifest_version": 1,
  "channel": "stable",
  "platform": "windows",
  "arch": "amd64",
  "version": "1.7.3",
  "release_seq": 42,
  "min_supported_version": "1.6.0",
  "url": "https://downloads.voltavpn.com/stable/windows/amd64/volta-1.7.3.exe",
  "sha256": "3f0e...64hex...",
  "created_at": "2026-02-12T10:30:00Z",
  "expires_at": "2026-02-19T10:30:00Z",
  "key_id": "prod-2026-01",
  "signature": "base64-ed25519-signature"
}
```

## Verification rules in client

1. Validate shape and scope:
   - `manifest_version == 1`
   - `channel/platform/arch` equal to client target
   - `url` is HTTPS and under `*.voltavpn.com`
2. Validate replay window:
   - `created_at <= now + skew`
   - `expires_at >= now - skew`
3. Validate anti-downgrade:
   - `release_seq > last_seen_release_seq`
   - `version > current_version`
4. Verify signature by `key_id`.
5. Download artifact and verify `sha256`.
6. Apply update only when all checks succeed (fail-closed).

## Downgrade and replay protection

- `release_seq` is monotonic and persisted on client after successful install.
- `version` check blocks same/older binary.
- `expires_at` blocks stale manifest replay.
- Channel and platform binding blocks cross-channel replay.

## Migration plan (no hard break)

### Phase 0: prepare infra

- Add update hosts (`updates.voltavpn.com`, `downloads.voltavpn.com`).
- Add release signing job in CI/CD.
- Keep private key in KMS/HSM or offline signing flow.

### Phase 1: dual mode clients

- New clients support signed manifest flow.
- Legacy version check remains as fallback for one transition period.
- Emit telemetry events for verification failures (without secrets).

### Phase 2: enforce secure updates

- Backend marks secure-manifest as required for client versions above a threshold.
- Disable legacy unsigned checks for new versions.
- Introduce mandatory update path with `min_supported_version`.

### Phase 3: hardening

- Key rotation drill each quarter.
- Incident runbook for compromised signing key.
- Optional threshold signing for high-risk channels.
