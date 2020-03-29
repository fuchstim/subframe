const { storage: config } = require('../../../config')
const { NamedError, errorTypes } = require('../../errors')
const { Logger } = require('../../logger')
const Database = require('../../utils/database')
const FileStorage = require('../../utils/file-storage')

class StorageHandler {
    constructor() {
        this.logger = new Logger('subframe/storage-handler')
        this.logger.debug(`Initializing StorageHandler with options ${JSON.stringify(config)}`)
        this.db = new Database({ path: config.path.database })
        this.fs = new FileStorage({ path: config.path.files, blockSize: config.blockSize, maxBlockCount: config.maxBlockCount })
        this.logger.debug('StorageHandler initialized.')
    }

    getFileInfo({ req, res, next }) {
        const { messageID } = req.params

        res.logger.debug(`Retrieving information for message '${messageID}'`)

        if (!messageID) {
            res.error = new NamedError(errorTypes.request.badRequest, 'No message ID supplied', { description: 'No message ID supplied' })
            return next()
        }
        if (!this._hasFile(messageID)) {
            res.error = new NamedError(errorTypes.request.unknownResource, 'Requesting unknown message')
            return next()
        }

        try {
            res.data = this.db.read(messageID)
        } catch (e) {
            res.error = e
        }

        next()
    }

    _hasFile(messageID) {
        return this.db.hasRecord(messageID)
    }

}

module.exports = new StorageHandler()