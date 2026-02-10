# VoltaVPN Client – High-Level Threat Model (Bootstrap Stage)

This is a **very early**, intentionally minimal threat model for the VoltaVPN client. It will evolve as features are added.

## Assets

- Future VPN configuration (servers, ports, protocols).
- Future authentication material (tokens, credentials, keys).
- User privacy metadata (which server is selected, connection status, logs).

Currently, none of these assets are implemented in code.

## Adversaries (High-Level)

- **Local attacker** with access to the same machine:
  - tries to read config and secrets from disk;
  - tries to manipulate the client binary or its configuration.
- **Network attacker** (once VPN logic exists):
  - tries to MITM client–server communication;
  - tries to inject malicious configuration or updates.
- **Malicious software** on the same host:
  - tries to exploit the client (GUI or core) to escalate privileges or exfiltrate data.

## Attack Surface (Current Skeleton)

- GUI surface:
  - Fyne-based desktop app;
  - local event handling and OS integrations.

What is **not** present yet:

- No network sockets opened by the client;
- No VPN protocol implementation;
- No crypto keys or credentials processed.

## High-Level Security Goals

As the client evolves, we aim for:

- **Confidentiality** of VPN credentials and keys (no accidental logs, no plaintext storage without strong justification).
- **Integrity** of configuration and binaries (defence against tampering).
- **Resilience** to common desktop app issues:
  - unsafe deserialization;
  - arbitrary code execution via external tools;
  - injection into shell/exec calls.

## Out of Scope (For Now)

- Protecting against a fully compromised host (e.g. active malware, rootkit).
- Protecting against hardware-level attacks.
- Providing formal proofs of protocol security (that will be handled at the protocol layer).

