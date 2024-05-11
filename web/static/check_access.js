function check_access(j, i, id){
  let key = prompt("Key:")
  let data = {
    "key" : key,
    "id" : id,
    "state": !document.getElementById(`active${j}${i}`).checked
  }

  fetch(`${window.origin}/check_status/access`, {
    method: 'POST',
    body: JSON.stringify(data),
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
        document.getElementById(`active${j}${i}`).checked = data["state"]
      }
    })
  })

}