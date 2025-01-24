<!--
SPDX-License-Identifier: GPL-3.0-only.
SPDX-FileCopyrightText: 2025 International Committee of the Red Cross.

This software is licensed under GPL v3.0 license.
See LICENSE file at the root of the project.
-->

# YALTF (Yet Another License Tool and Framework)

YALTF is a tool designed to scan and gather software licenses across multiple remote systems. It connects to target systems via SSH and collects information about installed packages and their associated licenses.

Developed by International Committee of the Red Cross (ICRC) in collaboration with Cortex Security S.A.

## Features

- Remote system scanning via SSH
- Parallel scanning of multiple targets
- Support for CentOS, Fedora, RHEL (Rocky) and OpenSUSE 
- Structured JSON output

## Prerequisites

- Go 1.22 or higher
- SSH (keybased) access to target systems

## Installation

```bash
# Clone the repository
git clone https://github.com/yaltf/yaltf.git

# Change to project directory
cd yaltf

# Build the project
go build
```

## Configuration

Create a configuration file named `config.toml` based on the provided template:

```bash
cp config.toml.template config.toml
```

Example configuration:

```toml
[common]
ssh_key_path = "$HOME/.ssh/id_rsa"

[targets.centos]
    host="192.168.1.10"
    port="22"
    user="admin"

[targets.rhel1]
    host="192.168.1.11"
    port="22"
    user="admin"  
```

## Usage

### Local Scanning

```bash
# Scan localhost
./yaltf scan -localhost
```

### Remote Scanning

```bash
# Scan all targets in config.toml
./yaltf scan
```

## Output

Results are stored in JSON format in the `results` directory. Each scan creates a new file with timestamp:

```json
{
  "version": "1.0",
  "scanned_at": "2025-01-23T07:12:35",
  "results": {
    "centos": {
      "package1": ["GPL-3.0"],
      "package2": ["MIT"]
    }
  }
}
```

## Development

### Supported OS
- CentOS: 8 to 10
- Fedora: 38 to 41
- RHEL: 9.4
- Rocky: 8.8
- OpenSUSE: 15.5

### Project Structure

```
.
├── main.go
├── config/
│   └── config.go
├── logging/
│   └── logging.go
├── models/
│   └── models.go
├── reporter/
│   └── reporter.go
├── scanner/
│   ├── scanner.go
│   └── parser.go
└── subcmds/
    └── scan.go
```
## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m -S 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Copyright and License

© 2025 International Committee of the Red Cross.

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.

The [International Committee of the Red Cross (ICRC)](https://www.icrc.org/) is an impartial, neutral and independent organization whose exclusively humanitarian mission is to protect the lives and dignity of victims of armed conflict and other situations of violence and to provide them with assistance.
The ICRC also endeavours to prevent suffering by promoting and strengthening humanitarian law and universal humanitarian principles.
Established in 1863, the ICRC is at the origin of the Geneva Conventions and the International Red Cross and Red Crescent Movement. It directs and coordinates the international activities conducted by the Movement in armed conflicts and other situations of violence.
