'use strict'

exports.slss = function (event, context, callback) {
  // Keep event loop rolling
  setTimeout(noop, Infinity)

  // TODO: call server binary
  callback(null, {})
}

function noop () {}
