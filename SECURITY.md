# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.2.x   | :white_check_mark: |
| < 0.2   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in ProxyDoctor, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

Instead, please email: [INSERT YOUR EMAIL HERE]

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### What to expect

- **Acknowledgement** within 48 hours
- **Assessment** within 1 week
- **Fix or mitigation** released as soon as possible

Contact: francomano on GitHub

## Security Considerations

ProxyDoctor handles sensitive data including:
- Proxy credentials (username/password in SOCKS5)
- TLS certificates and keys
- Network configuration data

All credential handling follows best practices:
- Credentials are never logged or stored on disk
- Proxy credentials are only held in memory during the connection lifecycle
- No telemetry or external data collection

## Scope

This security policy applies to:
- The CLI tool (`proxydoctor diagnose`)
- The HTTP server (`proxydoctor server`)
- The core engine and all check implementations
