from pathlib import Path
import re
import zipfile

root = Path(__file__).resolve().parents[1]
app = root / "remote_home_assistant"
configuration = (app / "config.yaml").read_text()
assert configuration.split("schema:\n", 1)[1].strip() == "pairing_token: password"
assert configuration.split("options:\n", 1)[1].split("schema:\n", 1)[0].strip() == 'pairing_token: ""'
version = re.search(r'^version: "([0-9.]+)"$', configuration, re.M).group(1)
for path in root.rglob("*"):
    if not path.is_file() or any(part in {".git", ".cache"} for part in path.relative_to(root).parts):
        continue
    if path.suffix in {".pem", ".key"} or path.name in {".env", "cloudflare", "known_hosts"}:
        raise SystemExit("Refusing private material in public repository")
    if path.suffix in {".go", ".py", ".yaml", ".md"}:
        data = path.read_text()
        if re.search(r"-----BEGIN (?:RSA |EC )?PRIVATE KEY-----", data):
            raise SystemExit("Refusing operator material in public repository")
output = root / ".cache" / f"remote-ha-app-{version}.zip"
output.parent.mkdir(exist_ok=True)
with zipfile.ZipFile(output, "w", zipfile.ZIP_DEFLATED) as archive:
    for path in sorted(app.rglob("*")):
        if path.is_file():
            archive.write(path, path.relative_to(root))
print(f"Validated token-only app {version}; archive: {output}")
