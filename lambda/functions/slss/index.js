'use strict'
const os = require('os')
const { spawn } = require('child_process')

const MAX_DELAY = 2147483647
const EVENT_REQUIRED_KEYS = ['port', 'method', 'password']

exports.handle = function (event, context, callback) {
  if (validateEvent(event)) print('event', event)
  else return callback(new Error(`Invalid event: ${JSON.stringify(event)}`))

  // Keep event loop rolling
  setTimeout(noop, MAX_DELAY)

  const ssOptions = ['-k', event.password, '-m', event.method, '-p', event.port]

  const server = spawn('./bin/shadowsocks_server', ssOptions)

  server.stdout.on('data', function (data) {
    print('ss_server stdout', data.toString())
    callback(null, { networkInterfaces: os.networkInterfaces() })
  })

  server.stderr.on('data', function (data) {
    print('ss_server stderr', data.toString())
  })

  server.on('close', function (code) {
    callback(new Error(`ss_server close, code: ${code}`))
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

function noop () {}
