import shutil
import os
import platform
import subprocess

home = os.environ["HOME"]

go_path = os.environ["GOPATH"]
if not os.path.exists(go_path):
    go_path = os.path.join(home, "go")
    os.mkdir(go_path)

go_bin = os.path.join(go_path, "bin")
if not os.path.exists(go_bin):
    os.mkdir(go_bin)

subprocess.run(["go", "build"])

bin = ""
if "Linux" in platform.platform():
    bin = "shutil"
else:
    bin = "shutil.exe"

dest = os.path.join(go_bin, bin)
shutil.move(bin, dest)
