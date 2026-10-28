# Remote HA for Home Assistant

Secure remote access to Home Assistant with one pairing token. This repository
contains only the public Home Assistant app and its agent. The hosted platform,
account management, PKI and operator infrastructure live in a separate private
repository. Building this app never requires that repository or its credentials.

## Install

Add this repository's HTTPS URL in **Settings → Apps → App store → Repositories**,
then install **Remote Home Assistant**. Create an installation at
[remote-ha.com](https://remote-ha.com), choose its address and generate a pairing
token. Paste the token into the app's only configuration field and start it.
Use your normal Home Assistant account at the installation address.

The app discovers Core through Supervisor. No public home IP, incoming router
port, user-provided certificate or Core restart is required. Existing paired
identity remains under `/data/rha` during updates; do not uninstall its data.

## Security model

The agent creates its device key locally. Private-CA mutual TLS protects the
outbound tunnel and short-lived signed authorization controls access. Cloudflare
and the platform terminate browser HTTPS and can process Home Assistant traffic.
Home Assistant authentication still applies. This is not end-to-end encryption
between your browser and Core. Core HTTPS, when configured, remains verified.

The app requests Supervisor's default role for Core discovery and migration of
its own settings. It has no Core administrator API permission, host networking,
privileged mode or Docker socket. Old configuration keys are removed through
the app's own Supervisor endpoint; the UI exposes only the pairing token.

## Develop and build

```sh
cd remote_home_assistant/src
go test -race ./...
```

```sh
docker buildx build --platform linux/amd64,linux/arm64 \
  -f remote_home_assistant/Dockerfile remote_home_assistant
```

`Dockerfile.agent` builds a standalone agent for controlled integration tests
or operator-managed Container deployments. Its endpoint overrides are not part
of the Home Assistant app's configuration. FRP downloads are version and SHA256
pinned. Release CI tests the agent and builds both supported architectures.

Report app bugs through this repository's Issues. Never include pairing tokens,
device keys, Supervisor tokens or unredacted diagnostics in a report. See
[SECURITY.md](SECURITY.md) for security reporting.
