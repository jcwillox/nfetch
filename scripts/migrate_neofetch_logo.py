import re
import sys

REGEX_COLOR_VAR = re.compile(r"\${\w(\d)}")

if len(sys.argv) < 2:
    print("please specify a file")
    exit(1)

path = sys.argv[1]

with open(path) as file:
    logo = file.read()

logo = REGEX_COLOR_VAR.sub(r"{{.C\1}}", logo)

# ensure logo ends with newline
if logo[-1] != "\n":
    logo += "\n"

print(logo)

with open(path, "w") as file:
    file.write(logo)
