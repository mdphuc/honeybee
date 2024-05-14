function addUser(form){
  let key = prompt("Key:")

  if (key == null){
    return
  }
  let fullname_col = prompt("Full Name col:")
  let email_col = prompt("Email col: ")

  let usr_col = prompt("Username col:")
  let pwd_col = prompt("Password col:")

  if (usr_col == null || pwd_col == null || fullname_col == null || email_col == null){
    return 
  }
  submit_form = document.getElementById(form)

  input = document.createElement('input')
  input.type = 'TEXT'
  input.name = 'key'
  input.value = key
  submit_form.appendChild(input)

  input1 = document.createElement('input')
  input1.type = 'TEXT'
  input1.name = 'usr_col'
  input1.value = usr_col
  submit_form.appendChild(input1)

  input2 = document.createElement('input')
  input2.type = 'TEXT'
  input2.name = 'pwd_col'
  input2.value = pwd_col
  submit_form.appendChild(input2)

  input3 = document.createElement('input')
  input3.type = 'TEXT'
  input3.name = 'fullname_col'
  input3.value = fullname_col
  submit_form.appendChild(input3)

  input4 = document.createElement('input')
  input4.type = 'TEXT'
  input4.name = 'email_col'
  input4.value = email_col
  submit_form.appendChild(input4)

  submit_form.submit()
  submit_form.removeChild(input)
  submit_form.removeChild(input1)
  submit_form.removeChild(input2)
  submit_form.removeChild(input3)
  submit_form.removeChild(input4)

}