'use strict'
const { spawn } = require('child_process')
const http = require('http')

const MAX_DELAY = 2147483647
const EVENT_REQUIRED_KEYS = [
  'port',
  'method',
  'password',
  'proxyHost',
  'proxyPort'
]

exports.handle = function (event, context, callback) {
  if (validateEvent(event)) print('event', event)
  else return callback(new Error(`Invalid event: ${JSON.stringify(event)}`))

  // Keep event loop rolling
  setTimeout(function () { callback(null) }, MAX_DELAY)

  const ssOptions = ['-k', event.password, '-m', event.method, '-p', event.port]
  addLogging(spawn('./bin/shadowsocks_server', ssOptions), 'ss_server')
  addLogging(spawn('./bin/ngrok', ['tcp', event.port]), 'ngrok')

  http.get(`http://${event.proxyHost}:${event.proxyPort}/`, function (res) {
    if (res.statusCode !== 200) {
      print('http_request', `bad status code error: ${res.statusCode}`)
    }
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
