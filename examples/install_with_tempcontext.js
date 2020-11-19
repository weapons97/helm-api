const fs = require('fs')
const http = require('http')

let { host, port } = process.env

let kubeconfig_path = process.argv[2]
let context = process.argv[3]
let namespace = process.argv[4]
let releaseName = process.argv[5]

let buff = fs.readFileSync(kubeconfig_path)
let kubeconfig = buff.toString()

let buff2 = fs.readFileSync(`install-values.yaml`)
let values = buff2.toString()

let data = {
  "releaseName": `${releaseName}`,
  "chartName": `redis`,
  "chartVersion": "",
  "namespace": `${namespace}`,
  "values": `${values}`,
  "repo": "testRepo",
  "dry_run": true,
  "kubeinfo": {
    kubeconfig: `${kubeconfig}`,
    context: `${context}`,
    namespace: `${namespace}`
  },
  "repoinfo": {
    entrys: [{
      name: `testRepo`,
      url: `https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts`
    }]
  }
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
