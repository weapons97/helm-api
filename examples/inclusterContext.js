const fs = require('fs')
const http = require('http')

let { host, port } = process.env

let namespace = process.argv[2]

let data = {
  "name": "test",
  "entry": {
    "name": "library",
    "url": "https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts",
    "username": "",
    "password": "",
    "certFile": "",
    "keyFile": "",
    "caFile": "",
    "insecure_skip_tls_verify": true
  },
  "kubeconfig": ``,
  "context": ``,
  "namespace": `${namespace}`,
  "incluster": true
}

let rbody = JSON.stringify(data)
let options = {
  hostname: host,
  port: port,
  path: "/v1/context",
  method: "POST",
  headers: {
    "Content-Length": Buffer.byteLength(rbody)
  }
}

console.log(`${host}:${port}`, "body:", rbody)

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
request.end(rbody)
