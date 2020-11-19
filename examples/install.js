const fs = require('fs')
const http = require('http')

let { host, port } = process.env

let releaseName = process.argv[2]
let namespace = process.argv[3]

let buff = fs.readFileSync(`install-values.yaml`)
let values = buff.toString()

let data = {
  "releaseName": `${releaseName}`,
  "chartName": `mariadb`,
  "chartVersion": "",
  "namespace": `${namespace}`,
  "values": `${values}`,
  "contextName": "test",
  "repo": "bitnami",
  "dry_run": false
}

let rbody = JSON.stringify(data)
let options = {
  hostname: host,
  port: port,
  path: "/v1/test/release",
  method: "POST",
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
