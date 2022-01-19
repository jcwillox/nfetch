import os
import shutil
from glob import glob

shutil.rmtree("dist/", ignore_errors=True)

for file in glob("nfetch.*.*"):
    os.remove(file)
