# Remote Home Assistant 1.0.1

Create an account and installation at https://remote-ha.com. Copy a pairing
token into **Pairing token**, the app's only setting, and start the app.
Open the assigned address and sign in with your Home Assistant account.

The app automatically discovers Core and establishes an outbound authenticated
mTLS tunnel. You do not need to generate certificates, expose a router port or
restart Core. Your device key and identity remain in `/data/rha` across upgrades.

When updating from older releases, the app removes obsolete settings through
its own Supervisor options endpoint. Saved pairing and `!secret` references are
preserved. Back up the app and its data before updating; keep the recovery key.

The app uses the default Supervisor role for Core information and its own
settings. It has no Core administrator token, host networking or Docker socket.
Existing Core HTTPS must validate against the standard trust store; certificate
verification is never disabled.

A local `/healthz` endpoint on port 8088 reports pairing and tunnel process
status without credentials. The platform dashboard reports recent activity.

Cloudflare and the platform terminate browser TLS and can process HA traffic.
The tunnel between the app and platform uses mutual TLS. Your Home Assistant
login remains independent from your platform account.
