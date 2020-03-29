const dateTime = require('node-datetime')
const outputOptions = require('./output-options')
const { logs: config } = require('../../config')

class Logger {
    constructor(context) {
        this.context = context
    }

    debug = (message) => { if (config.debug) { console.log(outputOptions.textColor.cyan, this.formatMessage('DBG', message)) } }
    info = (message) => console.log(outputOptions.textColor.green, this.formatMessage('INF', message))
    warn = (message) => console.log(outputOptions.textColor.yellow, this.formatMessage('WRN', message))
    error = (message) => console.log(outputOptions.textColor.red, this.formatMessage('ERR', message))
    fatal = (message) => console.log(outputOptions.textBackgroundColor.red, outputOptions.textColor.white, this.formatMessage('FAT', message))

    formatMessage = (status, message) => `[${status}] ${dateTime.create().format('Y-m-d H:M:S')} (${this.context}): ${message}`
}

function attachLoggerToRequest(req, res, next) {
    const requestLogger = new Logger(`request/${req.method.toLowerCase()}`)
    requestLogger.info(`'${req.path}'`)

    requestLogger.context = `${req.method.toLowerCase()}${req.path}`
    res.logger = requestLogger

    next()
}

module.exports = {
    Logger,
    attachLoggerToRequest
}