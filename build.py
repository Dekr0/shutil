import shutil
import os
import subprocess

subprocess.run(["go", "build"])
shutil.move("shutil.exe", os.path.join(os.environ["GOPATH"], "bin/shutil.exe"))
