import shutil
import os
import platform
import subprocess

# Prefer pathlib? Slightly slower, in exchange for abstraction.
from pathlib import Path

home_path = Path.home()

# os.environ['suchanenv'] WILL raise an exception if it does not exist
# os.environ.get('suchanenv') DOES NOT raise an exception if it does not exist

# Recall and consider the differences between `go env` and `env`
# At installation, a default GOPATH is set in `go env` but not `env`

# Assume go is installed, and some environment variable exists
bin_dir = os.environ.get("GOBIN")
go_dir = os.environ.get("GOPATH")
if bin_dir:
    install_dir = Path(bin_dir)
elif go_dir:
    install_dir = Path(go_dir) / "bin"
    install_dir.mkdir(exist_ok=True)
else:
    install_dir = home_path / "go" / "bin"
    # https://docs.python.org/3/library/pathlib.html#pathlib.Path.mkdir
    # Creates the directories /go/bin iff does not already exist
    install_dir.mkdir(parents=True, exist_ok=True)


subprocess.run(["go", "build"])

bin = ""
if "Linux" in platform.platform():
    bin = "shutil"
else:
    bin = "shutil.exe"

install_dest = install_dir / bin
shutil.move(bin, install_dest)
