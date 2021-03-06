const http = require('http')

let { host, port } = process.env

let contextName = process.argv[2]

let options = {
  hostname: host,
  port: port,
  path: `/v1/${contextName}/chart`,
  method: "GET",
  headers: {
  }
}

var returned = ""
var request = new http.ClientRequest(options, (res) => {
  res.on('data', function (chunk) {
    returned += chunk
  })
  res.on('end', () => {
    let beauty = JSON.stringify(JSON.parse(returned), null, 4)
    console.log('Response: \n', beauty)
  })
}).on('error', function (e) {
  console.log('Problem with request: ' + e)
})
request.end()
