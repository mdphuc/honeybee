
let imageName = ""

function expose_port(ports, image_name){
  console.log("yo")
  imageName = image_name
  console.log(imageName)
  if (ports !== "null"){
    if (document.getElementById(`ports_div_${imageName}`) != null){
      document.getElementById("port_divv").removeChild(document.getElementById(`ports_div_${imageName}`))
    }else{
      ports_div = document.createElement('div')
      ports_div.id = `ports_div_${imageName}`
      ports_div.innerHTML += `<form id="expose_ports_${imageName}" action="/expose_port" method="post" style="width: 100%; margin-left: 40px"></form>`
      ports_div.innerHTML += `<div id="ports_list_${imageName}"></div>`
      ports_div.innerHTML += `<button onclick="add_port('${imageName}')">Add</button><br><br>`
      ports_div.innerHTML += `<button onclick="changeport('${imageName}')">Submit</button>`
      document.getElementById("port_divv").appendChild(ports_div)
  
      ports_list = document.getElementById(`ports_list_${imageName}`)
      ports_div.style.display = "block"
      ports_array = ports.split(",")
  
      for(let i = 0; i < ports_array.length; i++){
        if (document.getElementById(`expose_port_${imageName}_${i}`) != null){
          document.getElementById(`expose_port_${imageName}_${i}`).remove()
        }
        ports_list.innerHTML += `<div id="expose_port_div_${imageName}_${i}"><input id="expose_port_${imageName}_${i}" value=${ports_array[i]}><button onclick="remove_port_list(${i})" style="margin-left:10px">X</button><br><br></div>`
      }
      
    }


  }
}

function add_port(imageName){
  ports_list = document.getElementById(`ports_list_${imageName}`)
  num_ports_list = ports_list.querySelectorAll("input").length
  document.getElementById(`ports_list_${imageName}`).innerHTML += `<div id="expose_port_div_${imageName}_${num_ports_list}"><input id="expose_port_${imageName}_${num_ports_list}"><button onclick="remove_port_list(${num_ports_list})" style="margin-left:10px">X</button><br><br></div>`
}

function remove_port_list(num_ports_list){
  document.getElementById(`ports_list_${imageName}`).removeChild(document.getElementById(`expose_port_div_${imageName}_${num_ports_list}`))
}

function changeport(imageName){
  let key = prompt("Key:")
  let d = ""

  valid_input = true

  ports_div = document.getElementById(`ports_list_${imageName}`)
  let div_split_port = ports_div.querySelectorAll("div")
  for (let i = 0; i < div_split_port.length; i++){
    ports = div_split_port[i].querySelector("input").value
    if (ports.value != "" && parseInt(ports.value) != NaN) {
      if (i != ports.length - 1){
        d += `${ports[i].value},`
      }else{
        d += `${ports[i].value}`
      }
    }else{
      valid_input = false
      break
    }
  }

  if (valid_input){
    ports_form = document.getElementById(`expose_ports_${imageName}`)

    input = document.createElement('input')
    input.type = 'TEXT'
    input.name = 'ports_to_update'
    input.value = d
    ports_form.appendChild(input)

    machine = document.createElement('input')
    machine.type = 'TEXT'
    machine.name = 'machine'
    machine.value = imageName
    ports_form.appendChild(machine)

    key_input = document.createElement('input')
    key_input.type = 'TEXT'
    key_input.name = 'key'
    key_input.value = key
    ports_form.appendChild(key_input)

    ports_form.submit()
    ports_form.removeChild(input)
    ports_form.removeChild(key_input)
    ports_form.removeChild(machine)
  }else{
    alert("Blank port field or port entered is not a number")
  }

}
