from web.main import app, db, User
import hashlib
import datetime
import os

# os.remove("database.db")
# db.create_all()
usr = User(username="yo", password=hashlib.sha256(b"yo").hexdigest(), role="user", key="",date_created=datetime.datetime.now())
db.session.add(usr)
db.session.commit()
