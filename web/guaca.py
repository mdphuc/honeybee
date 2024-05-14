import requests
import random
import string
 
USERNAME = "guacadmin"
PASSWORD = "guacadmin"

def get_token():
  url = "http://127.0.0.1:8080/api/tokens"

  headers = {
    'Content-Type' : 'application/x-www-form-urlencoded',
  }

  data = { 'username' : USERNAME, 'password' : PASSWORD }

  response = requests.post(url, headers=headers, data=data)
  
  return response.json()['authToken']

# token = get_token()

def add_user(user, token):
  url = "http://127.0.0.1:8080/api/session/data/postgresql/users?token={0}".format(token)

  data = {
  "username": user["username"],
  "password": user["fullname"] + ''.join(random.choices(string.ascii_uppercase + string.digits, k=10)),
  "attributes": {
      "disabled": "",
      "expired": "",
      "access-window-start": "",
      "access-window-end": "",
      "valid-from": "",
      "valid-until": "",
      "timezone": None,
      "guac-full-name": user["fullname"],
      "guac-email-address": user["email"],
      "guac-organizational-role": "User"
    } 
  }

  response1 = requests.post(url, json=data)

  url = "http://127.0.0.1:8080/api/session/data/postgresql/users/{0}/permissions?token={1}".format(user["username"], token)

  data = [
    {
      "op": "add",
      "path": "/systemPermissions",
      "value": "CREATE_CONNECTION"
    },
    {
      "op": "add",
      "path": "/systemPermissions",
      "value": "CREATE_CONNECTION_GROUP"
    }
  ]

  response2 = requests.patch(url, json=data)

  if (response1.status_code == 200 and response2.status_code == 200):
    print("Succeed!")

def add_admin(admin, token):
  url = "http://127.0.0.1:8080/api/session/data/postgresql/users?token={0}".format(token)

  data = {
  "username": admin["username"],
  "password": admin["fullname"] + ''.join(random.choices(string.ascii_uppercase + string.digits, k=10)),
  "attributes": {
      "disabled": "",
      "expired": "",
      "access-window-start": "",
      "access-window-end": "",
      "valid-from": "",
      "valid-until": "",
      "timezone": None,
      "guac-full-name": admin["fullname"],
      "guac-email-address": admin["email"],
      "guac-organizational-role": "Admin"
    } 
  }

  response1 = requests.post(url, json=data)

  url = "http://127.0.0.1:8080/api/session/data/postgresql/users/{0}/permissions?token={1}".format(admin["username"], token)

  data = [
    {
      "op": "add",
      "path": "/systemPermissions",
      "value": "ADMINISTER"
    }
  ]

  response2 = requests.patch(url, json=data)

  if (response1.status_code == 200 and response2.status_code == 200):
    print("Succeed!")


# token = get_token()
# print(token)

