from flask import Flask, render_template, request, abort, redirect, url_for, session, jsonify, make_response
from flask_sqlalchemy import SQLAlchemy
from flask_login import UserMixin, login_required, LoginManager, current_user, login_user, logout_user
from datetime import  timedelta, datetime
from itsdangerous import URLSafeTimedSerializer
from sqlalchemy.sql import func 
import os
import hashlib
import random
import subprocess 
import datetime
import logg
from convert_time import convert_time
import random
from graph import *
from werkzeug.utils import secure_filename
import pandas as pd
import string
import asym
import guaca

#BEAVER##DUH

number = random.randint(10**5,10**15)

# admin_key = os.environ.get("HoneyBeeAdminID") 

keywords = ["Port", "AllowTerminalOnWeb", "Exec", "PasswordProtected", "AllowAccess", "IP"]

PORT_RANGE = [i for i in range(49153, 65536)]
ALLOWED_EXTENSIONS = {'csv'}

def config():
  f = open("../config/default.conf")
  flag = {}
  data = f.readlines()
  for d in data:
    keyword = d.split("=")[0]
    if d.split("=")[0] in keywords:
      flag[keyword] = d.split("=")[1].split("\n")[0]
  return flag

flag = config()

for keyword in flag.keys():
  if keyword == "AllowAccess" and flag[keyword] != "":
    ip_allow = str(flag[keyword]).split()
  else:
    ip_allow = ["127.0.0.1"]
  if keyword == "Port":
    PORT = int(str(flag[keyword]).split()[0])
  if keyword == "AllowTerminalOnWeb":
    AllowTerminalOnWeb = flag[keyword]
  if keyword == "Exec":
    EXEC = flag[keyword]
  if keyword == "PasswordProtected":
    PasswordProtected = flag[keyword]
  if keyword == "IP":
    IP = flag[keyword]
  else:
    IP = '127.0.0.1'

def allowed_file(filename):
  return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS

admin_key = hashlib.sha256(b"dawg").hexdigest()

app = Flask(__name__)
app.config['SECRET_KEY'] = hashlib.sha256(b"app").hexdigest()
app.config["SQLALCHEMY_DATABASE_URI"]= "sqlite:///database.db"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
app.permanent_session_lifetime = timedelta(days=2)
app.config['suppress_callback_exception'] = True

login_manager = LoginManager()
login_manager.session_protection = "strong"
login_manager.login_view = "login"
login_manager.init_app(app)

@login_manager.user_loader
def load_user(user_id):
  return User.query.get(int(user_id))

db = SQLAlchemy(app)

s = URLSafeTimedSerializer(app.secret_key)


class User(db.Model, UserMixin):
  __tablename__ = "user"
  id = db.Column(db.Integer, primary_key=True, unique=True)
  fullname = db.Column(db.String(30), nullable=False, unique=True)
  email = db.Column(db.String(40), nullable=False, unique=True)
  username = db.Column(db.String(20), nullable=False, unique=True)
  password = db.Column(db.String(50), nullable=False)
  active = db.Column(db.Boolean, default=True, nullable=False)
  role = db.Column(db.String(10), nullable=False)
  date_created = db.Column(db.String(20), nullable=False)
  last_login = db.Column(db.String(20), nullable=False)
  key = db.Column(db.String(30), default="null")
  key_user = db.Column(db.String(30), default="null")
  guac_password = db.Column(db.String(50), default="")

  def __init__(self, username, password, role, key, date_created, last_login, active, key_user, fullname, email, guac_password):
    super().__init__()
    self.username = username
    self.password = password
    self.role = role
    self.key = key
    self.key_user = key_user
    self.date_created = date_created
    self.last_login = last_login
    self.active = active
    self.fullname = fullname
    self.email = email
    self.guac_password = guac_password
  
  def __repr__(self):
      return f"User('{self.username}','{self.password}','{self.role}','{self.key}','{self.date_created}','{self.last_login}','{self.active}','{self.key_user}','{self.fullname}','{self.email}','{self.guac_password}')"

@app.before_request
def block_method():
  if request.remote_addr not in ip_allow:
    logg.logging("main.py Block IP", None)
    abort(403)

@app.route("/403", methods=["GET"])
def forbidden():
  logg.logging("main.py unauthorized access", current_user.username)
  abort(403)

@app.route("/<type>", methods=["GET", "POST"])
@login_required
def home(type):
  if request.method == "GET":
    if current_user.is_authenticated:
      logs = log_time_calculate()
      if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
        return render_template("admin.html", logs = logs, name = current_user.username)
      else:
        if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
          if type != "graph":
            log_remove = []
            for i in range(len(logs)):
              if logs[i].split()[4] != current_user.username:
                log_remove.append(i)
            for lr in log_remove[::-1]:
              print(lr)
              logs.pop(lr)
            return render_template("normal.html", logs = logs, name = current_user.fullname)
          else:
            di = get_status("normal")
            return render_template("normalgraph.html", data = di, graph="true")
        else:
          logg.logging("main.py unauthorized access", current_user.username)
          abort(403)       
    else:
      return render_template("login.html")
  else:
    if current_user.is_authenticated:
      if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
        info = check_access(type.split("_")[0], "all")
        return render_template("admin.html",  logs = info[0], headers = info[1], data = info[2], access = info[3], name = current_user.username)
      else:
        key = request.form.get("key")
        if key == asym.decrypt(current_user.key_user, current_user.username):
          logs = log_time_calculate()
          log_remove = []
          for i in range(len(logs)):
            if logs[i].split()[4] != current_user.username:
              log_remove.append(i)
          for lr in log_remove[::-1]:
            print(lr)
            logs.pop(lr)
          # for lr in log_remove:
          #   logs.pop(lr)
          di = get_status("normal")
          print(di)
          headers = ["ID", "STATUS", "REPOSITORY", "TAG", "IMAGE ID", "SIZE", "CREATED AT", "CONTAINER ID", "UP FOR", "PORTS"]
          
          if len(type.split("_")) == 1:
            return render_template("normal.html", headers = headers, data = di, logs = logs, access="true",graph="false", name = current_user.fullname)
          else:
            if type.split("_")[1] == "graph":
              graph = Graph()
              graph.get_access_point()
              graph.get_node(di)
              graph.generate()
              return render_template("normal.html", headers = headers, data = di, logs = logs, access="true", graph="true", name = current_user.fullname)
        else:
          logg.logging("main.py unauthorized access", current_user.username)
          abort(403)
    else:
      return render_template("login.html")

@app.route("/guac_setup/<state>", methods=["POST"])
@login_required
def guac_setup(state):
  if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    req = request.get_json()
    if req["key"] == asym.decrypt(current_user.key, current_user.username):
      if state == "initiate":
        original = os.getcwd()
        os.chdir("./guac")
        dockercompose = subprocess.Popen(["docker", "compose", "up", "-d"], stdout=subprocess.PIPE)
        (target_container, err) = dockercompose.communicate()
        os.chdir(original)
      elif state == "user_setup":
        users = User.query.all()
        usrs = []
        for usr in users:
          if usr.username != "admin" and usr.id != 1 and asym.decrypt(usr.guac_password, usr.username) == "":
            usrs.append({"fullname" : usr.fullname, "username" : usr.username, "email" : usr.email, "role" : usr.role})
        try:
          token = guaca.get_token()
        except:
          logg.logging("main.py guacamole container down or error within container", current_user.username)
          return redirect(url_for("home", type="index"))
        try:
          for u in usrs:
            if usr.role == "admin":
              guaca.add_admin(usr)
            elif usr.role == "user":
              guaca.add_user(usr)
        except:
          logg.logging("main.py error when add user", current_user.username)
          return redirect(url_for("home", type="index"))
         
    else:
      logg.logging("main.py unauthorized access", current_user.username)
      abort(403)
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)



@app.route("/login", methods=["GET", "POST"])
def login():
  if request.method == "POST":
    session.permanent = True
    username = request.form.get("uname")
    password = request.form.get("password")

    user = User.query.filter_by(username = username, password = hashlib.sha256(password.encode()).hexdigest(), active = True).first()
    if user is not None:
      login_user(user)
      user.last_login = datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S")
      db.session.commit()
      return redirect(url_for("home", type="index"))
    else:
      return render_template("login.html", message = "Failed")
  if request.method == "GET":
    return render_template("login.html", message = "")

@app.route("/logout", methods=["POST"])
@login_required
def logout():
  try:
    session.pop('username')
    session.pop('password')
  except:
    pass
  logout_user()
  return redirect(url_for("home", type="index"))

@app.route("/getKey", methods=["POST"])
@login_required
def getKey():
  if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    req = request.get_json()
    if hashlib.sha256(req["key"].encode()).hexdigest() == current_user.password:
      res = make_response(jsonify({"key": asym.decrypt(current_user.key, current_user.username)}), 200)
      return res
    else:
      logg.logging("main.py unauthorized access", current_user.username)
      abort(403)  
  elif current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    req = request.get_json()
    if hashlib.sha256(req["key"].encode()).hexdigest() == current_user.password:
      if req["target"] == "KEY":
        res = make_response(jsonify({"key": asym.decrypt(current_user.key_user, current_user.username)}), 200)
        return res
      elif req["target"] == "GUACKEY":
        gp = asym.decrypt(current_user.guac_password, current_user.username)
        if gp == "":
          res = make_response(jsonify({"key": "None"}), 200)
          return res
        else:
          res = make_response(jsonify({"key": gp}), 200)
          return res
         
    else:
      logg.logging("main.py unauthorized access", current_user.username)
      abort(403) 
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)  
 

@app.route("/register", methods=["POST"])
@login_required
def register():
  if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    key = request.form.get("key")
    usr_col = request.form.get("usr_col")
    pwd_col = request.form.get("pwd_col")
    fullname_col = request.form.get("fullname_col")
    email_col = request.form.get("email_col")
    if key == asym.decrypt(current_user.key, current_user.username):
      file = request.files['user_file']
      if file.filename != "":
        if file and allowed_file(file.filename):
          filename = secure_filename(file.filename)
          file.save('./{0}'.format(filename))
          try:
            data = pd.read_csv("./{0}".format(filename), usecols=[fullname_col, email_col, usr_col, pwd_col])
            for i in range(len(data[usr_col])):
              key_user = asym.encrypt(hashlib.sha256(data[fullname_col][i] + ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))).hexdigest())
              usr = User(username=data[usr_col][i], password=hashlib.sha256(data[pwd_col][i].encode()).hexdigest(), role="user", key="null", key_user=asym.encrypt(key_user, current_user.username), date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True, fullname=data[fullname_col][i], email=data[email_col][i],guac_password=asym.encrypt("", current_user.username))
              db.session.add(usr)
              db.session.commit()
          except:
            logg.logging("main.py error when reading csv file", current_user.username)
      return redirect(url_for("home", type="index"))
    else:
      logg.logging("main.py unauthorized access", current_user.username)
      abort(403)  
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)  

@app.route("/reg", methods=["POST"])
@login_required
def reg():
  if current_user.is_authenticated and current_user.get_id() == "1" and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    fullname = request.form.get("name")
    email = request.form.get("email")
    username = request.form.get("username")
    password = request.form.get("password")

    if fullname == "" or email == "" or username == "" or password == "":
      logg.logging("main.py missing information when trying to add user", current_user.username)
      
    try:
      key = asym.encrypt(hashlib.sha256(fullname + ''.join(random.choices(string.ascii_uppercase + string.digits, k=10))).hexdigest())
      usr = User(username=username, password=password, role="admin", key=key, key_user="null", date_created=datetime.datetime.utcnow().strftime("%Y/%m/%d %H:%M:%S"), last_login="null", active=True, fullname=fullname, email=email)
      db.session.add(usr)
      db.session.commit()
    except:
      logg.logging("main.py error when adding user", current_user.username)
    
    return redirect(url_for("home", type="index"))
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)  


@app.route("/dockerfile", methods=["POST"])
@login_required
def dockerfile():
  if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    return redirect(url_for("home", type="index"))
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)


@app.route("/nodockerfile", methods=["POST"])
@login_required
def nodockerfile():
  if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    req = request.get_json()
    if req['key'] == asym.decrypt(current_user.key_user, current_user.username):
      text = request.form["dockerfile-edit"]
      Dockerfile_initial = open("../Dockerfile", "w")
      Dockerfile_initial.write(text)
      Dockerfile_initial.close()
      Dockerfile_default = open("../Dockerfile", "r")
      lines = Dockerfile_default.readlines()

      Dockerfile = open("../Dockerfile2", "w")
      place = 0
      count_break = 0
      for i in range(len(lines)):
        if lines[i] == "\n" and place == 0 and count_break == 0:
          place = i
          j = i+1
          while True:
            if j < len(lines):
              if lines[j] != "\n":
                Dockerfile.writelines("\n"*(count_break-1))
                break
              else:
                count_break += 1
              j += 1
            else:
              break
        if i < place + count_break + 1 and i > 0:
          pass
        else:
          place = 0
          count_break = 0
          Dockerfile.writelines(lines[i])

      Dockerfile_default.close()
      Dockerfile.close()

      os.remove("../Dockerfile")
      os.rename("../Dockerfile2", "../Dockerfile")
      return redirect(url_for("home", type="index"))
    else:
      logg.logging("main.py unauthorized access", current_user.username)
      abort(403)
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)

def check_access(type, target):
  if current_user.is_authenticated and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    key = request.form["key"]
    user = User.query.filter_by(role = "admin").first()
    if key == asym.decrypt(user.key, current_user.username):
      if type == "docker":
        di = get_status(target)
        logs = log_time_calculate()
        headers = ["ID", "STATUS", "REPOSITORY", "TAG", "IMAGE ID", "SIZE", "CREATED AT", "CONTAINER ID", "UP FOR", "PORTS"]
        return (logs, headers, di, "true")
      elif type == "user":
        logs = log_time_calculate()

        di = []
        users = User.query.all()
        for user in users:
          if user.username != "admin" and user.key != admin_key:
            if asym.decrypt(user.guac_password, user.username) == "":
              di.append([user.id, user.username, user.fullname, user.email, user.active, user.role, user.date_created, user.last_login, "-"])
            else:
              di.append([user.id, user.username, user.fullname, user.email, user.active, user.role, user.date_created, user.last_login, "UP"])     
        headers = ["ID", "USERNAME", "FULL NAME", "EMAIL", "ACTIVE", "ROLE", "DATE CREATED", "LAST LOGIN", "GUAC"]
        return (logs, headers, di, "true")
    else:
      return redirect(url_for("home", type="index"))
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)

@app.route("/check_status/<type>", methods=["POST"])
@login_required
def check_status(type):
  if current_user.is_authenticated and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    req = request.get_json()
    if type != "access":
      di = get_status("all") 
      res = make_response(jsonify({"data": di, "length_depth" : len(di[0])}), 200)
    else:
      if asym.encrypt(req['key']) == current_user.key:    
        res = make_response(jsonify({"access": "true"}), 200)
        user = User.query.filter_by(id = int(req['id'])).first()
        if req["state"] == True:
          user.active = False
        elif req["state"] == False:
          user.active = True
        db.session.commit()
      else:
        res = make_response(jsonify({"access": "false", "state": req['state']}), 200)
        logg.logging("main.py unauthorized access", current_user.username)
      # res = make_response(jsonify({"data": "yo"}), 200)
      # print (req)
    return res
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)

@app.route("/check_status_user", methods=["POST"])
@login_required
def check_status_user(type):
  if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    req = request.get_json()
    if type != "access":
      di = get_status("all") 
      res = make_response(jsonify({"data": di, "length_depth" : len(di[0])}), 200)
    else:
      if req['key'] == asym.decrypt(current_user.key, current_user.username):    
        res = make_response(jsonify({"access": "true"}), 200)
        user = User.query.filter_by(id = int(req['id'])).first()
        if req["state"] == True:
          user.active = False
        elif req["state"] == False:
          user.active = True
        db.session.commit()
      else:
        res = make_response(jsonify({"access": "false", "state": req['state']}), 200)
        logg.logging("main.py unauthorized access", current_user.username)
      # res = make_response(jsonify({"data": "yo"}), 200)
      # print (req)
    return res
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)



def get_status(target):
  proc = subprocess.Popen(["docker", "images", "--format", "{{.Repository}} {{.Tag}} {{.ID}} {{.Size}} {{.CreatedAt}} "], stdout=subprocess.PIPE)
  (docker_images, err) = proc.communicate()

  docker_images = docker_images.decode().split("\n")

  di = []
  ID = 1
  for docker_image in docker_images:
    if len(docker_image) > 2:
      general = docker_image.split()
      name = general[0]
      if target == "all":
        users = User.query.all()

        names = [user.username for user in users]
        temp = name.split("-")
        if len(temp) < 2:
          exist = False
        else:
          exist = name.split("-")[1] in names
      else:
        temp = name.split("-")
        if len(temp) < 2:
          exist = False
        else:
          if name.split("-")[1] == current_user.username:
            exist = True
          else:
            exist = False   
      if len(name.split("-")) == 3 and name.split("-")[0] == "honeybee" and name.split("-")[1] == current_user.username and exist:
        tag = general[1]
        id = general[2]
        size = general[3]
        created_at = ' '.join(general[4::])
        read = subprocess.Popen(["docker", "ps", "--format", "{{.Image}}~{{.ID}}~{{.Status}}~{{.Ports}}"], stdout=subprocess.PIPE)
        (target_image, err) = read.communicate()
        target_image = target_image.decode().split("\n")
        print(target_image)
        if len(target_image) == 1:
          di.append([ID, "DOWN", name, tag, id, size, created_at, "null", "null", "null"])
          ID += 1
        else:
          exist_further = False
          for ti in target_image:
            if name in ti:
              ports = ti.split("~")[3].split(">")[1].split("-")
              for i in range(len(ports)):
                ports[i] = ports[i].split("/")[0].strip()
              print(ports)
              di.append([ID, "UP", name, tag, id, size, created_at, ti.split("~")[1], ti.split("~")[2], ','.join(ports)])
              exist_further = True
              break
          if not exist_further:
            di.append([ID, "DOWN", name, tag, id, size, created_at, "null", "null", "null"])
          ID += 1
  
  return di

@app.route('/expose_port', methods = ["POST"])
@login_required
def expose_port():
  if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    key = request.form["key"]
    if key == asym.decrypt(current_user.key_user, current_user.username):
      ports_to_update = "".join(request.form["ports_to_update"].split())
      print(ports_to_update)
      if ports_to_update != "":
        image_name = request.form["machine"]
        di = get_status("normal")
        exist = False
        for d in di:
          if image_name in d:
            exist = True
            break

        if exist:
          read = subprocess.Popen(["docker", "ps", "--format", "{{.Ports}}"], stdout=subprocess.PIPE)
          (target_image, err) = read.communicate()
          target_image = target_image.decode().split("\n")

          can_set_port_host = []
          occupied_port_host_list = []
          print(target_image)

          target_image.remove(target_image[len(target_image) - 1])

          for ti in target_image:
            occupied_port_host = ti.split("->")[0].split(":")[1].split("-")
            for oph in occupied_port_host:
              occupied_port_host_list.append(oph)

          for ophl in occupied_port_host_list:
            for p in PORT_RANGE:
              if p != int(ophl) and p not in can_set_port_host:
                can_set_port_host.append(p)
                break
          
          print(can_set_port_host)

          proc = subprocess.Popen(["docker", "commit", image_name, image_name], stdout=subprocess.PIPE)
          (docker_commit, err) = proc.communicate()

          docker_update_port = ["docker", "run", "--name", image_name]

          proc = subprocess.Popen(["docker", "stop", image_name], stdout=subprocess.PIPE)
          (docker_stop, err) = proc.communicate()

          proc = subprocess.Popen(["docker", "system", "prune", "-f"], stdout=subprocess.PIPE)
          (docker_prune, err) = proc.communicate()

          # for p in can_set_port_host:
          previous_set = 0
          for ptd in ports_to_update:
            p_index = random.randint(0, len(can_set_port_host) - 1)
            while p_index == previous_set:
              p_index = random.randint(0, len(can_set_port_host) - 1)
            previous_set = p_index

            p = can_set_port_host[p_index]
            docker_update_port.append("-p")
            docker_update_port.append("{0}:{1}".format(p, ptd))

          docker_update_port.append("-itd")
          docker_update_port.append(image_name)
          docker_update_port.append("bash")
          # print(docker_update_port)
          proc = subprocess.Popen(docker_update_port, stdout=subprocess.PIPE)
          (docker_update_port, err) = proc.communicate()

        else:
          logg.logging("main.py unauthorized access", current_user.username)
          abort(403)
      return redirect(url_for("home", type="index"))
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)

@app.route("/update", methods = ["POST"])
@login_required
def update():
  if current_user.is_authenticated and asym.decrypt(current_user.key, current_user.username) == admin_key and current_user.key_user == "null":
    key = request.form.get("key")
    usr_col = request.form.get("usr_col")
    pwd_col = request.form.get("pwd_col")
    fullname_col = request.form.get("fullname_col")
    if key == asym.decrypt(current_user.key, current_user.username):
      file = request.files['user_file']
      if file.filename != "":
        if file and allowed_file(file.filename):
          filename = secure_filename(file.filename)
          file.save('./{0}'.format(filename))
          try:
            data = pd.read_csv("./{0}".format(filename), usecols=[fullname_col, usr_col, pwd_col])
            for i in range(len(data[usr_col])):
              usr = User(fullname=data[fullname_col][i])
              if usr == None:
                logg.logging("main.py error when reading csv file", current_user.username)
              else:
                if usr.username != data[usr_col][i]:
                  usr.username = data[usr_col][i]
                  db.session.commit()
                if usr.password != hashlib.sha256(data[pwd_col][i].encode()).hexdigest():
                  usr.password = hashlib.sha256(data[pwd_col][i].encode()).hexdigest()
                  db.session.commit()
          except:
            logg.logging("main.py error when reading csv file", current_user.username)
      return redirect(url_for("home", type="index"))

  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)


@app.route("/check_iframe", methods=["POST"])
@login_required
def check_iframe():
  if current_user.is_authenticated and current_user.key == "null" and asym.decrypt(current_user.key_user, current_user.username) != "null":
    req = request.get_json()
    if req['key'] == asym.decrypt(current_user.key_user, current_user.username):
      read = subprocess.Popen(["docker", "ps", "--format", "{{.Image}}~{{.ID}}~{{.Status}}~{{.Ports}}"], stdout=subprocess.PIPE)
      (target_image, err) = read.communicate()
      target_image = target_image.decode().split("\n")
      port = 0
      for ti in target_image:
        if ti != '':
          ti = ti.split("~")
          if ti[0] == req['image']:
            ti_tail = ti[3].split("->")[1].split("/")[0].split("-")
            ti_head = ti[3].split("->")[0].split(":")[1].split("-")
            print(ti_head)
            for i in range(len(ti_tail)):
              if str(ti_tail[i]) == str(2005):
                port = ti_head[i]
      print(port)
      res = make_response(jsonify({"access": "true", "port": port}), 200)
    else:
      res = make_response(jsonify({"access": "false"}), 200)
    return res
  else:
    logg.logging("main.py unauthorized access", current_user.username)
    abort(403)
  
def check_log():
  with open('../log/app.log', 'r') as f:
    lines = f.readlines()

  if "\n" not in lines[len(lines)-1]:
    with open('../log/app.log', 'a') as f:
      f.write("\n")

def log_time_calculate():
  http_launch_time = ""
  log_list = []

  check_log()

  with open("../log/app.log", "r") as f:
    log_times = f.readlines()
    for log_time in log_times:
      year = log_time.split()[0].split("/")[0]
      month = log_time.split()[0].split("/")[1]
      day = log_time.split()[0].split("/")[2]
      hour = log_time.split()[1].split(":")[0]
      minute = log_time.split()[1].split(":")[1]
      second = log_time.split()[1].split(":")[2]

      if log_time.split()[4] == "httplaunch":
        http_launch_time = year + month + day + hour + minute + second

    for log_time in log_times:
      year = log_time.split()[0].split("/")[0]
      month = log_time.split()[0].split("/")[1]
      day = log_time.split()[0].split("/")[2]
      hour = log_time.split()[1].split(":")[0]
      minute = log_time.split()[1].split(":")[1]
      second = log_time.split()[1].split(":")[2]

      log_time_ = year + month + day + hour + minute + second

      if log_time_ > http_launch_time and log_time.split()[4] != "httplaunch":
        log_list.append(log_time)
        
  return log_list

if __name__ == "__main__":
  db.create_all()
  app.run("127.0.0.1", 5500, debug=True)
