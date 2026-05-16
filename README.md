# FleetSSL cPanel

Automatic, free SSL certificates for cPanel &amp; WHM, powered by [Let's Encrypt](https://letsencrypt.org).

FleetSSL cPanel issues and **automatically renews** trusted SSL certificates for every
domain on a cPanel/WHM server. It provides:

- A **WHM interface** for server administrators
- A **cPanel interface** for end users (Jupiter &amp; Paper Lantern themes)
- A **background service** that renews certificates before they expire
- An **AutoSSL provider**, so Let's Encrypt can be used as cPanel's AutoSSL backend
- `HTTP-01` and `DNS-01` domain validation, with automatic fallback between them
- A `le-cp` command-line interface

This is a community-maintained fork of the FleetSSL cPanel plugin, which was
open-sourced by FleetSSL.com under the MIT license. No warranty or technical
support is provided — see [LICENSE.md](LICENSE.md).

## Requirements

- A working **cPanel &amp; WHM** installation
- A **64-bit (x86_64)** server
- **root** shell access
- A supported operating system:

  | Package | Operating systems |
  | ------- | ----------------- |
  | `.rpm`  | CentOS 7, AlmaLinux 8/9, Rocky Linux 8/9, CloudLinux 7+ (any RHEL-family &ge; 7) |
  | `.deb`  | Ubuntu 20.04 LTS |

## Install (RHEL / AlmaLinux / Rocky / CloudLinux)

Log in to the server as **root** and run:

```bash
yum install -y https://github.com/persianopencart/fleetssl-cpanel-new/releases/latest/download/letsencrypt-cpanel.x86_64.rpm
```

On AlmaLinux / Rocky / CloudLinux 8 and newer, `dnf` works exactly the same way.

> Use `yum`/`dnf` rather than `rpm -i`. The package depends on `bc`, and only
> `yum`/`dnf` will resolve and install that dependency for you.

That single command performs the **entire** installation:

1. Downloads the package and installs it to `/opt/fleetssl-cpanel`
2. Checks that the server is running cPanel on a supported OS
3. Stops any previously running copy of the service
4. Writes the default configuration to `/etc/letsencrypt-cpanel.conf`
5. Installs and starts the `letsencrypt-cpanel` background service
6. Registers the WHM plugin, the cPanel plugin, and the AutoSSL provider
7. Registers the service with cPanel's `chkservd` monitoring
8. Adds the Apache configuration required for HTTP domain validation

When it finishes you will see `--- Installation complete ---`.

### Self-signed hostname certificate

The plugin talks to the cPanel/WHM APIs over HTTPS. If the server's hostname
still has a self-signed certificate, the installer cannot verify it and enables
`insecure` mode in `/etc/letsencrypt-cpanel.conf` so the plugin can still
operate. For a fully secure setup, install a valid certificate on the server
hostname first.

### Newer or unrecognised operating systems

The installer refuses to run on operating systems it does not recognise. To
force installation anyway, prefix the command with `FLEETSSL_SKIP_OS_CHECK=y`:

```bash
FLEETSSL_SKIP_OS_CHECK=y yum install -y https://github.com/persianopencart/fleetssl-cpanel-new/releases/latest/download/letsencrypt-cpanel.x86_64.rpm
```

## Install (Ubuntu)

```bash
apt-get install -y https://github.com/persianopencart/fleetssl-cpanel-new/releases/latest/download/letsencrypt-cpanel.amd64.deb
```

## After installing

**Server administrators (WHM)** — log in to WHM and open the **FleetSSL cPanel**
plugin (search the WHM menu for "FleetSSL" or "Let's Encrypt"). From there you
can configure global behaviour, issue a certificate for the server hostname, and
review renewals.

**End users (cPanel)** — a **"Let's Encrypt&trade; SSL"** icon appears in cPanel
under the **Security** section. If users cannot see it, enable the feature for
their package in **WHM &rarr; Feature Manager**.

**Command line** — the CLI is installed as `le-cp`:

```bash
le-cp self-test     # check that every part of the plugin works
le-cp help          # list all commands
```

**Service management:**

```bash
systemctl status letsencrypt-cpanel      # or: service letsencrypt-cpanel status
systemctl restart letsencrypt-cpanel
journalctl -u letsencrypt-cpanel -f      # follow the logs
```

**Important paths:**

| Path | Purpose |
| ---- | ------- |
| `/opt/fleetssl-cpanel/` | Installed program files |
| `/etc/letsencrypt-cpanel.conf` | Configuration |
| `/var/lib/letsencrypt-cpanel.db` | Plugin database |
| `/usr/local/bin/le-cp` | CLI symlink |

## Verify the installation

```bash
le-cp self-test
```

This confirms the configuration is readable, the server can reach Let's Encrypt,
the cPanel/WHM APIs respond, and the background service is running.

## Upgrading

Re-run the install command — `yum`/`dnf` (or `apt-get`) upgrades the package in
place and the service is restarted automatically:

```bash
yum install -y https://github.com/persianopencart/fleetssl-cpanel-new/releases/latest/download/letsencrypt-cpanel.x86_64.rpm
```

## Installing a specific version

Every release publishes the same asset filenames. To pin a version, use its
release tag instead of `latest`:

```bash
yum install -y https://github.com/persianopencart/fleetssl-cpanel-new/releases/download/<tag>/letsencrypt-cpanel.x86_64.rpm
```

Browse available versions on the
[Releases page](https://github.com/persianopencart/fleetssl-cpanel-new/releases).

## Uninstalling

```bash
yum remove letsencrypt-cpanel        # RHEL family
apt-get remove letsencrypt-cpanel    # Ubuntu
```

Uninstalling removes the plugin, the service and the WHM/cPanel integration.
Per-account certificate data stored in each user's cPanel `nvdata` is **kept**,
as are `/etc/letsencrypt-cpanel.conf` and `/etc/letsencrypt-cpanel.licence` —
delete those by hand if you want a completely clean system.

## Building from source

All Go dependencies are vendored, so the build needs no network access.
Build tools required: Go 1.19+, `fpm`, `rpm`/`rpmbuild`, `gawk`, `dos2unix`.

```bash
git clone https://github.com/persianopencart/fleetssl-cpanel-new.git
cd fleetssl-cpanel-new

make            # build just the binary (letsencrypt.live.cgi)
make rpm        # build the versioned .rpm and .deb packages
make dist       # build the packages and copy them to release asset names
```

## Releases

Releases are automated. Every push to `main` runs
[semantic-release](https://semantic-release.gitbook.io/): the version is derived
from the [Conventional Commits](https://www.conventionalcommits.org/) in the
history, a GitHub release is created, and the `.rpm` and `.deb` packages are
built and attached to it automatically.

## License

MIT — see [LICENSE.md](LICENSE.md). Original work &copy; FleetSSL.com.
