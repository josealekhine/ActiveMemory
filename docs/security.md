---
#   /    Context:                     https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Security
icon: lucide/shield
---

## Reporting Vulnerabilities

At `ctx` we take security very seriously.

If you discover a security vulnerability in `ctx`, please report it responsibly.

**Do NOT open a public issue for security vulnerabilities.**

### Email

Send details to **security@ctx.ist**

### GitHub Private Reporting

1. Go to the [Security tab](https://github.com/ActiveMemory/ctx/security)
2. Click "Report a vulnerability"
3. Provide a detailed description

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

| Stage              | Timeframe                              |
|--------------------|----------------------------------------|
| Acknowledgment     | Within 48 hours                        |
| Initial assessment | Within 7 days                          |
| Resolution target  | Within 30 days (depending on severity) |

## Security Design

`ctx` is designed with security in mind:

- **No secrets in context**: The constitution explicitly forbids storing
  secrets, tokens, API keys, or credentials in `.context/` files
- **Local only**: ctx runs entirely locally with no external network calls
- **No code execution**: ctx reads and writes Markdown files only; it does 
  not execute arbitrary code
- **Git-tracked**: All context files are meant to be committed, so they should
  never contain sensitive data

## Best Practices

1. **Review before committing**: Always review `.context/` files before committing
2. **Use .gitignore**: If you must store sensitive notes locally,
   add them to `.gitignore`
3. **Drift detection**: Run `ctx drift` to check for potential issues

## Attribution

We appreciate responsible disclosure and will acknowledge security researchers
who report valid vulnerabilities (*unless they prefer to remain anonymous*).
