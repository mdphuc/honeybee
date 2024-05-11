from main import app, db, User
import hashlib
import datetime
import os

# os.remove("database.db")
# db.create_all()
# usr = User(username="admin", password=hashlib.sha256(b"admin").hexdigest(), role="admin", key="yo",date_created=datetime.datetime.now())
# db.session.add(usr)
# db.session.commit()


users = User.query.all()
for user in users:
  print(user.active)

for user in users:
  print(user)