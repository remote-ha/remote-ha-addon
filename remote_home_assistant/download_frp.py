import hashlib
import io
import os
from pathlib import Path
import tarfile
import urllib.request

VERSION = "0.71.1"
CHECKSUMS = {
    "amd64": "1d9875210972a5d1a62b2dc4c48df45ab0af05d3ca1b5e494a014428259df803",
    "arm64": "8edf590b0c396fa42f9c147121a76af83cfc2002b08f5aef0ea1008c43482019",
}
arch = os.environ["TARGETARCH"]
data = urllib.request.urlopen(f"https://github.com/fatedier/frp/releases/download/v{VERSION}/frp_{VERSION}_linux_{arch}.tar.gz", timeout=120).read()
if hashlib.sha256(data).hexdigest() != CHECKSUMS[arch]:
    raise RuntimeError("FRP checksum mismatch")
archive = tarfile.open(fileobj=io.BytesIO(data))
for name in ("frpc",):
    path = Path("/usr/local/bin") / name
    path.write_bytes(archive.extractfile(f"frp_{VERSION}_linux_{arch}/{name}").read())
    path.chmod(0o755)
