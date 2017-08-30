'use strict'
const { execFile } = require('child_process')

const MAX_DELAY = 2147483647
const EVENT_REQUIRED_KEYS = ['port', 'method', 'password']

exports.slss = function (event, context, callback) {
  if (validateEvent(event)) printEvent(event)
  else return callback(new Error(`Invalid event: ${JSON.stringify(event)}`))

  // Keep event loop rolling
  setTimeout(noop, MAX_DELAY)

  const ssOptions = [
    `-k ${event.password}`,
    `-m ${event.method}`,
    `-p ${event.port}`
  ]

  execFile('./bin/shadowsocks_server', ssOptions, function (err, stdout, stderr) {
    if (err) return callback(err)

    callback(null, event)
  })
}

function validateEvent (event) {
  for (const key of EVENT_REQUIRED_KEYS) if (!event[key]) return false
  return true
}

function printEvent (event) {
  console.log('--------------- event ---------------')
  console.log(JSON.stringify(event))
  console.log('-------------------------------------')
}

function noop () {}
