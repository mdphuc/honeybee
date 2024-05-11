from main import app, db, User
import hashlib
import datetime
import os

try:
  os.remove("database.db")
except:
  pass
db.create_all()
usr1 = User(username="admin", password=hashlib.sha256(b"admin").hexdigest(), role="admin", key="yo", key_user="null", date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True)
usr2 = User(username="yo", password=hashlib.sha256(b"yo").hexdigest(), role="user", key="null", key_user="damn", date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True)
db.session.add(usr1)
db.session.add(usr2)
db.session.commit()
