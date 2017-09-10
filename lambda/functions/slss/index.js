'use strict'
const { spawn } = require('child_process')

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

  const ssOptions = [
    '-k', event.password,
    '-m', event.method,
    '-p', event.port
  ]
  const server = spawn('./bin/shadowsocks_server', ssOptions)

  const gostOptions = [
    `-D`,
    `-L=:${event.port}`,
    `-F=tcp://${event.proxyHost}:${event.proxyPort}`
  ]
  const gost = spawn('./bin/gost', gostOptions)

  server.stdout.on('data', (data) => print('ss_server stdout', data.toString()))
  server.stderr.on('data', (data) => print('ss_server stderr', data.toString()))
  server.on('close', (code) => callback(new Error(`ss_server close, code: ${code}`)))

  gost.stdout.on('data', (data) => print('gost stdout', data.toString()))
  gost.stderr.on('data', (data) => print('gost stderr', data.toString()))
  gost.on('close', (code) => callback(new Error(`gost close, code: ${code}`)))
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
