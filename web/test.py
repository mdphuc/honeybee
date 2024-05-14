import os
import subprocess

o = os.getcwd()
print(os.getcwd())
os.chdir("./guac")
print(os.getcwd())
os.chdir(o)
print(os.getcwd())

# dockercompose = subprocess.Popen(["docker", "compose", "up", "-d"], stdout=subprocess.PIPE)
# (target_container, err) = dockercompose.communicate()
