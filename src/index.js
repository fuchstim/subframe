const express = require('express')
const app = express()

const { express: config } = require('../config')
const coordinatorNode = require('./services/coordinatorNode')
const storageNode = require('./services/storageNode')
const { Logger, attachLoggerToRequest } = require('./logger')
const formatResult = require('./utils/format-result')
const { version } = require('../package.json')

const logger = new Logger('subframe/index')

logger.info(`Starting SuBFraMe Version ${version}...`)

app.use(attachLoggerToRequest)

app.use('/coordinator', coordinatorNode.router)
app.use('/storage', storageNode.router)

app.use(formatResult)

app.listen(config.port, config.host)

logger.info(`Started SuBFraMe Server. Running on ${config.host}:${config.port}`)