# Security Policy

This document outlines our security policy and procedures for reporting and handling security vulnerabilities in Gwid Core. We take security seriously and appreciate your efforts to responsibly disclose any issues you find.

## Supported Versions

We are committed to providing security updates for the following versions of Gwid Core:

* **Current Stable Version** (`v0.1.0a`) - All security patches will be backported to this version.
* **Previous Stable Version** (`v0.0.1a`) - Critical security patches will be considered for backporting.

Older versions are not officially supported and may not receive security updates. We strongly recommend upgrading to a supported version.

## Reporting a Vulnerability

If you discover a security vulnerability in Gwid Core, please do **NOT** open a public GitHub issue. Instead, we ask you to report it responsibly via our dedicated security channel.

**Preferred Reporting Method:**

* **Email:** Send an email to `info@gwid.io`.

**Please include the following information in your report:**

1.  **Vulnerability Description:** A clear and concise description of the vulnerability.
2.  **Steps to Reproduce:** Detailed steps to reliably reproduce the vulnerability. This is crucial for us to confirm and fix the issue.
    * Include code snippets, configuration, or commands if applicable.
    * Specify the version of Gwid Core and Go you were using.
    * Detail your operating system and environment.
3.  **Impact:** Explain the potential impact of the vulnerability (e.g., data loss, unauthorized access, denial of service).
4.  **Proof of Concept (Optional but Recommended):** If possible, provide a small, self-contained proof-of-concept (PoC) that demonstrates the vulnerability.
5.  **Proposed Fix (Optional):** If you have a suggestion for a fix, please include it.

**Response Time:**

* We will acknowledge your report within **2 business days**.
* We will provide an estimated timeline for a fix and public disclosure (if applicable) within **7 business days** after initial acknowledgement, depending on the complexity of the issue.

## Disclosure Policy

Our general policy is to follow a responsible disclosure model. This means:

1.  **Private Communication:** We will communicate with you privately throughout the remediation process.
2.  **Fix Development:** We will work to develop a fix for the vulnerability.
3.  **Coordinated Disclosure:** Once a fix is ready, we will coordinate with you on a public disclosure date. This typically involves:
    * Releasing a new version of Gwid Core with the fix.
    * Publishing a security advisory (e.g., GitHub Security Advisory, CVE).
    * Attributing the discovery to you (unless you prefer to remain anonymous).

We aim for a reasonable disclosure timeline, typically within **60 days** (e.g., 30-90 days) from the initial report, allowing users sufficient time to update.

## Best Practices for Secure Development (for Contributors and Maintainers)

As an open-source Go project, we encourage and strive for secure coding practices. Here are some guidelines:

* **Input Validation and Sanitization:** Always validate and sanitize all user inputs to prevent injection attacks (e.g., SQL injection, XSS).
* **Error Handling:** Handle errors gracefully and avoid exposing sensitive information in error messages.
* **Secure Dependencies:**
    * Regularly update Go modules to their latest stable versions.
    * Utilize tools like `govulncheck` to scan dependencies for known vulnerabilities.
    * Prefer Go's standard library `crypto` packages for cryptographic operations over third-party alternatives unless absolutely necessary and thoroughly vetted.
* **Logging and Monitoring:** Implement comprehensive logging to detect and investigate suspicious activities.
* **Avoid `unsafe` Package:** Use of the `unsafe` package should be strictly avoided unless there's a compelling reason and the code has undergone rigorous security review.
* **Code Review:** All code changes undergo thorough peer review to identify potential security flaws.
* **Static Analysis:** Utilize static analysis tools like `gosec` as part of the CI/CD pipeline to automatically detect common security weaknesses.
* **Minimize Attack Surface:** Design applications with the principle of least privilege and expose only necessary functionalities.

## Public Security Advisories

All public security advisories for Gwid Core will be published on:

* GitHub Security Advisories tab, `https://github.com/livepeer-gwid/gwid-core/security/advisories`

Thank you for helping us keep Gwid Core secure!
