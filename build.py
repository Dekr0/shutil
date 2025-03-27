import shutil
import os
import platform
import subprocess

subprocess.run(["go", "build"])

if "Linux" in platform.platform():
    shutil.move("shutil", os.path.join(os.environ["GOPATH"], "bin/shutil"))
else:
    shutil.move("shutil.exe", os.path.join(os.environ["GOPATH"], "bin/shutil.exe"))
