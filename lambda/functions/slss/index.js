'use strict'
const { spawn } = require('child_process')
const http = require('http')

const MAX_DELAY = 2147483647
const EVENT_REQUIRED_KEYS = [
  'port',
  'method',
  'password',
  'proxyURL',
  'ngrokToken'
]

exports.handle = function (event, context, callback) {
  if (validateEvent(event)) print('event', event)
  else return callback(new Error(`Invalid event: ${JSON.stringify(event)}`))

  // Keep event loop rolling
  setTimeout(() => callback(null), MAX_DELAY)

  const ssOptions = ['-k', event.password, '-m', event.method, '-p', event.port]
  addLogging(spawn('./bin/shadowsocks_server', ssOptions), 'ss_server')

  getNgrokAddress(event.port, event.ngrokToken)
    .then((addr) => {
      http.get(`http://${event.proxyURL}/?ss_server_addr=${addr}`, ({ statusCode }) => {
        if (statusCode === 200) return print('http_request', 'success')
        print('http_request', `bad status code error: ${statusCode}`)
      })
    })
    .catch((error) => print('ngrok_error', error))
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

function getNgrokAddress (port, token) {
  return new Promise(function (resolve, reject) {
    const ngrok = spawn('./bin/ngrok', [
      'tcp',
      port,
      '-log=stdout',
      '--log-level=debug',
      '--region=au',
      `--authtoken=${token}`
    ])

    ngrok.stdout.on('data', function (data) {
      const dataString = data.toString()
      if (!dataString.includes('tcp://')) return

      const i = dataString.lastIndexOf('tcp://')
      return resolve(dataString.slice(i + 'tcp://'.length, i + dataString.slice(i).indexOf(' ')))
    })
    ngrok.stderr.on('data', reject)
    ngrok.on('close', reject)
  })
}
