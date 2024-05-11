function sleep(ms){
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function get_status(){
  while (true){
    let entry = {
      text : "check"
    }
    console.log("run")

    fetch(`${window.origin}/check_status/status`, {
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
        for (let i = 0; i < data["data"].length; i++) {
          for (let j = 0; j < data["length_depth"]; j++){
            if (j == 1){
              if (data["data"][i][1] == "UP"){
                document.getElementById(`${i}${j}`).innerHTML = '<span style="color: #32CD32">●</span>'
              }else{
                document.getElementById(`${i}${j}`).innerHTML = '<span style="color: red">●</span>'                 
              } 
            }else{
              document.getElementById(`${i}${j}`).innerHTML = data["data"][i][j]
            }
          }
        } 
        
      })
    })
    await sleep(100*1000)
  }
}

