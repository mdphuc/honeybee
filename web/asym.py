import base64
import os
from cryptography.fernet import Fernet
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
import hashlib

admin_key = hashlib.sha256(b"dawg").hexdigest()

def encrypt(msg, salt):
  password = admin_key.encode()
  s = salt.encode()
  kdf = PBKDF2HMAC(
      algorithm=hashes.SHA256(),
      length=32,
      salt=s,
      iterations=480000,
  )

  key = base64.urlsafe_b64encode(kdf.derive(password))
  f = Fernet(key)
  return f.encrypt(msg.encode()).decode()


def decrypt(enc, salt):
  password = admin_key.encode()
  s = salt.encode()
  kdf = PBKDF2HMAC(
      algorithm=hashes.SHA256(),
      length=32,
      salt=s,
      iterations=480000,
  )

  key = base64.urlsafe_b64encode(kdf.derive(password))
  f = Fernet(key)

  return f.decrypt(enc.encode()).decode()

# a = encrypt("", "10")
# print(a)
# print(decrypt(a, "10"))

# print(decrypt("gAAAAABmQdZZURUJplnCGDAcT4LVMKVWzzZx-eQ5qTr8wASBdWq2gvWDivV8RZh4lbzy4rIJHEOu_8pJnbWwihjkzHH7TSoslQ=="))
# print(f.decrypt(b"gAAAAABmQdmoLkJ-vx8EzFo7-Bk0fiBXo_Fz-VBSsUh-_AdARZkAqJo9BqNvJ8wMgq5q6MOrJJEwnqGxJwU8pBM160fwTyh6lQ==").decode())
# print(decrypt("gAAAAABmQdPgpM5LziRnueCcm-Mc6-q5ia4q6x0FHYS64w9rjoboerD7PTNogTUpeWOVudisB8dYPI7iCl4WAzj8UWKzfTj50A=="))