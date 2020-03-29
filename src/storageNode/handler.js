const { storage: config } = require('../../config')
const { NamedError, errorTypes } = require('../errors')
const { Logger } = require('../logger')

class StorageHandler {
    constructor() {
        this.logger = new Logger('subframe/storage-handler')
        this.logger.debug(`Initializing StorageHandler with options ${JSON.stringify(config)}`)
        this.path = config.path
        this.maxSize = config.maxSize
        this.logger.debug('StorageHandler initialized.')
    }

    getFileInfo({ req, res, next }) {
        const { messageID } = req.params
        if (!messageID) {
            res.error = new NamedError(errorTypes.request.badRequest, 'No message ID supplied', { description: 'No message ID supplied' })
            return next()
        }
        if (!this.hasFile(messageID)) {
            res.error = new NamedError(errorTypes.request.unknownResource, 'Requesting unknown message')
            return next()
        }

        next()
    }

    hasFile(messageID) {
        return false
    }

}

module.exports = new StorageHandler()