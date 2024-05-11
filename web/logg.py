import subprocess

def logging(message, user):
  if user != "admin" and user != None:
    message = user + " " + message

  proc = subprocess.Popen(["go", "run", "./root.go", message], stdout=subprocess.PIPE)
  (m, err) = proc.communicate()


