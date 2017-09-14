'use strict'
const { spawn } = require('child_process')
const http = require('http')

const MAX_DELAY = 2147483647
const EVENT_REQUIRED_KEYS = [
  'port',
  'method',
  'password',
  'proxyHost',
  'proxyPort',
  'ngrokToken'
]

exports.handle = function (event, context, callback) {
  if (validateEvent(event)) print('event', event)
  else return callback(new Error(`Invalid event: ${JSON.stringify(event)}`))

  // Keep event loop rolling
  setTimeout(function () { callback(null) }, MAX_DELAY)

  addLogging(spawn('./bin/ngrok', ['authtoken', event.ngrokToken]), 'ngrok_auth')

  const ssOptions = ['-k', event.password, '-m', event.method, '-p', event.port]
  addLogging(spawn('./bin/shadowsocks_server', ssOptions), 'ss_server')

  getNgrokAddress(addLogging(spawn('./bin/ngrok', ['tcp', event.port]), 'ngrok'))
    .then(function (addr) {
      http.get(`http://${event.proxyHost}:${event.proxyPort}/?ss_server_addr${addr}`, function (res) {
        if (res.statusCode !== 200) {
          print('http_request', `bad status code error: ${res.statusCode}`)
        }
        print('http_request', 'success')
      })
    })
}

function validateEvent (event) {
  for (const key of EVENT_REQUIRED_KEYS) if (!event[key]) return false
  return true
}

function print (name, event) {
  console.log(`--------------- ${name} ---------------`)
  console.log(JSON.stringify(event))
  console.log('---------------------------------------')
}

function addLogging (emitter, name) {
  emitter.stdout.on('data', (data) => print(`${name} stdout`, data.toString()))
  emitter.stderr.on('data', (data) => print(`${name} stderr`, data.toString()))
  emitter.on('close', (code) => print(`${name} close`, code))

  return emitter
}

function getNgrokAddress (ngrok) {
  return new Promise(function (resolve, reject) {
    ngrok.on('data', function (data) {
      const dataString = data.toString()
      if (!dataString.includes('tcp://')) return

      const i = dataString.lastIndexOf('tcp://')

      return resolve(dataString.slice(i + 'tcp://'.length, i + dataString.slice(i).indexOf(' ')))
    })

    ngrok.on('error', reject)
  })
}
