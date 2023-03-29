import os
from pprint import pprint
import time
import dotenv
import sys

pprint(os.environ.get("test_var"))
pprint(os.environ.get("node_env"))

while True:
    print("Hello World", flush=True)
    time.sleep(1)