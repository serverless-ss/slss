'use strict'
const { exec } = require('child_process')

exports.handler = function (event, context, callback) {
  printObjectInfo(event, context)

  // Try to let this function run as long as possible
  setTimeout(function () { callback(null) }, Infinity)

  exec('', function (err, stdout, stderr) {
    if (err) callback(err)
  })
}

function printObjectInfo (...objects) {
  for (let object in objects) console.info(JSON.stringify(object))
}
