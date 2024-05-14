from main import app, db, User
import hashlib
import datetime
import os
import asym

admin_key = hashlib.sha256(b"dawg").hexdigest()
 
try:
  os.remove("database.db")
except:
  pass
db.create_all()
usr1 = User(username="admin", password=hashlib.sha256(b"admin").hexdigest(), role="admin", key=asym.encrypt(admin_key, "admin"), key_user="null", date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True, fullname="null",email="null",guac_password=asym.encrypt("guacadmin", "admin"))
usr2 = User(username="yo", password=hashlib.sha256(b"yo").hexdigest(), role="user", key="null", key_user=asym.encrypt("damn", "yo"), date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True, fullname="Phuc Mai", email="mdphuc@gmail.com",guac_password=asym.encrypt("", "yo"))
db.session.add(usr1)
db.session.add(usr2)
# usr = User(username="guacadmin", password="guacadmin", role="admin", key="yo", key_user="null", date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True)
db.session.commit()
