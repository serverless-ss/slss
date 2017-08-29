'use strict'
const { execFile } = require('child_process')

const DELAY = 2147483647

exports.slss = function (event, context, callback) {
  // Keep event loop rolling
  setTimeout(noop, DELAY)

  execFile('shadowsocks-server', [], function (err, stdout, stderr) {
    if (err) return callback(err)
  })

  // TODO: call server binary
  callback(null, {})
}

function noop () {}
