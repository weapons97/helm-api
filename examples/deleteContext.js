const http = require('http')

let { host, port } = process.env

let name = process.argv[2]

let options = {
  hostname: host,
  port: port,
  path: `/v1/context/${name}`,
  method: "DELETE",
  headers: {
  }
}

var returned = ""
var request = new http.ClientRequest(options, (res) => {
  res.on('data', function (chunk) {
    returned += chunk
  })
  res.on('end', () => {
    console.log('Response: ', returned)
  })
}).on('error', function (e) {
  console.log('Problem with request: ' + e)
})
request.end()
