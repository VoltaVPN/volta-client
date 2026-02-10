# Security Policy for VoltaVPN Client

This document describes how to report vulnerabilities and how we approach security for the VoltaVPN client.

## Reporting a Vulnerability

- Please **do not** create public GitHub issues for security problems.
- Instead, contact the maintainers privately (contact details will be added once the project is public and maintainers are finalized).
- Provide as much detail as possible:
  - affected version / commit hash;
  - environment (OS version, architecture);
  - steps to reproduce;
  - potential impact (confidentiality / integrity / availability).

We aim to:

- Acknowledge reports within a reasonable time frame;
- Provide a rough timeline for investigation and fixes, when possible;
- Credit you in the release notes if you want, once a fix is shipped.

## Scope

This repository currently contains:

- A Go-based GUI skeleton (Fyne) for a future VPN client;
- No real VPN functionality or cryptographic protocol logic yet.

As we add actual VPN and crypto code, this policy will be expanded with:

- more specific guarantees;
- stricter rules for code review;
- more detailed threat models and hardening guidelines.

## Principles

- **No security theater**: we do not claim that the client is “unhackable” or “military grade”.
- **Secure by design** where feasible:
  - minimal use of unsafe operations;
  - clear separation of GUI, core logic, and OS-specific code;
  - minimal privileges required to run.
- **Defence in depth**:
  - static analysis and code review for security-sensitive changes;
  - separation of secrets from code and configuration.

## What We Do Not Promise

- We cannot guarantee that the software is free of vulnerabilities.
- We do not currently provide an SLA for security fixes.
- We do not currently run a bug bounty program.

