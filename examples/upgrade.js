const fs = require('fs')
const http = require('http')

let { host, port } = process.env

let releaseName = process.argv[2]

let buff = fs.readFileSync(`upgrade-values.yaml`)
let values = buff.toString()

let data = {
  "contextName": "test",
  "releaseName": `${releaseName}`,
  "chartName": `mariadb`,
  "chartVersion": "",
  "namespace": "default",
  "values": `${values}`,
  "repo": "bitnami",
  "dry_run": false,
  "history_max": 0,
  "reset_values": false
}


let rbody = JSON.stringify(data)
let options = {
  hostname: host,
  port: port,
  path: "/v1/test/release",
  method: "PUT",
  headers: {
    "Content-Length": Buffer.byteLength(rbody)
  }
}

console.log("body:", rbody)

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
request.end(rbody)
