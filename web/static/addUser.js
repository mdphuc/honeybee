function addUser(form){
  let key = prompt("Key:")
  let usr_col = prompt("Username col:")
  let pwd_col = prompt("Password col:")
  if (key == null || usr_col == null || pwd_col == null){
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

  submit_form.submit()
  submit_form.removeChild(input)
  submit_form.removeChild(input1)
  submit_form.removeChild(input2)

}