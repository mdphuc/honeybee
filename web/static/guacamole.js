function guacamole(state){
  let key = prompt("Key:")
  let entry = {
    "key": key
  }
  fetch(`${window.origin}/guac_setup/${state}`, {
    method: 'POST',
    body: JSON.stringify(entry),
    cache: "no-cache",
    headers: new Headers({
      "content-type" : "application/json"
    })
  }).then(function(response){
    if (response.status !== 200){
      console.log(`Failed : ${response.status}`)
      return
    }

    response.json().then(function(data){
      if (data["access"] == "true"){
        document.getElementById("message_").innerHTML = "<p style='color: green; font-weight: bold'>Succeed</p>"
      }else if (data["access"] == "false"){
        document.getElementById("message_").innerHTML = "<p style='color: red; font-weight: bold'>Failed</p>"
      }
    })
  })

}