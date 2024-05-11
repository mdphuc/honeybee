from web.main import app, PORT, IP, PasswordProtected

if __name__ == "__main__":
  app.run(IP, PORT, debug=True)
  
