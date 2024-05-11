function ttyd(dockeri){
  let key = prompt("Key: ")
  let data = {
    "key": key,
    "image": dockeri
  }
  if (document.getElementById(`iframe_ttyd_${dockeri}`) != null){
    document.getElementById(`iframe_ttyd_${dockeri}`).remove
  }

  fetch(`${window.origin}/check_iframe`, {
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
        div_ttyd = document.getElementById("ttyd")
        div_ttyd.innerHTML = ""
        iframe = document.createElement('iframe')
        iframe.src = `http://${window.location.hostname}:${data['port']}`
        console.log(iframe.src)
        iframe.id = `iframe_ttyd_${dockeri}`
        iframe.style = 'width:100%; height:400px; display:block'
        div_ttyd.append(iframe)
        document.getElementById("ttyd").style.display = 'block'
        document.getElementById("ttyd_close").style.display = 'block'
        document.getElementById("ttyd").innerHTML += "<br>"        
      }else if (data["access"] == "false"){
        console.log("Failed")
        document.getElementById("message_").innerHTML += "<p style='color:red; font-weight:bold'>Access denied!</p>"
      }
    })
  })

}