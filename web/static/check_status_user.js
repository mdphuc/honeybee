function check_status_user(form){
  let key = prompt("Key:")
  if (key == null){
    return
  }
  submit_form = document.getElementById(form)
  input = document.createElement('input')
  input.type = 'TEXT'
  input.name = 'key'
  input.value = key
  submit_form.appendChild(input)
  submit_form.submit()
  submit_form.removeChild(input)
}