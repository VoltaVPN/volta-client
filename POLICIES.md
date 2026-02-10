# VoltaVPN Client – Secure Development Policies

These policies apply to all code in this repository. They are intentionally strict to reduce the risk of common security issues.

## 1. No Logging of Secrets

**Policy:** It is forbidden to log any sensitive data, including but not limited to:

- private keys (any form);
- access tokens, refresh tokens, auth codes;
- usernames + passwords;
- session identifiers, cookies;
- raw VPN configuration that includes credentials or keys;
- full request/response bodies that may contain secrets.

**Guidelines:**

- Prefer structured logging with explicit whitelisting of fields.
- When in doubt, **do not log** the value; log a hash, a boolean flag, or a redacted placeholder instead.
- Logs must be safe to share for debugging without exposing secrets.

## 2. No Secrets in the Repository

**Policy:** No private keys or long-term secrets may be stored in:

- the Git repository (tracked or untracked);
- commit history;
- example configs, unless they are obviously fake/dummy keys.

This includes, but is not limited to:

- TLS private keys;
- VPN server private keys;
- API keys and tokens for infrastructure or third parties.

**Guidelines:**

- Use environment variables, OS keychains, or secure secret managers for any real secrets.
- Example configs must either:
  - not contain secret-like values; or
  - contain clearly dummy values and be labeled as such.
- If a secret is accidentally committed:
  - rotate it as soon as possible;
  - remove it from the repository history if feasible;
  - document the incident in an internal security log.

## 3. Restriction on exec/shell Usage

**Policy:** Using OS-level command execution is **disallowed by default**.

This includes:

- `os/exec` in Go;
- invoking shells like `cmd.exe`, `powershell`, `sh`, `bash` via any mechanism;
- embedding scripting engines that execute untrusted or semi-trusted code.

**Allowed only with:**

- a strong, documented justification;
- prior review by a maintainer with security experience;
- a clear threat analysis for the concrete use case.

If such a use case ever appears, it must:

- never pass unsanitized user input to the shell;
- avoid constructs that permit arbitrary command chaining;
- be covered by tests that validate input sanitization and safe behavior.

## 4. General Secure Coding Practices

- Prefer immutable data structures and clear ownership where possible.
- Avoid global state for security-relevant data.
- Use standard library crypto and well-reviewed libraries instead of home-grown primitives.
- Treat any user- or network-provided data as untrusted and validate it strictly.

## 5. Code Review Expectations

- Any change that touches:
  - crypto;
  - key management;
  - authentication or authorization;
  - IPC, networking, or OS integration;
  - logging of potentially sensitive data;
  must go through at least one additional reviewer.

