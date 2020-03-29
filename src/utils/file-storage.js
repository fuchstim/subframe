const fs = require('fs')
const path = require('path')

const { Logger } = require('../logger')
const { NamedError, errorTypes } = require('../errors')

class FileStorage {
    constructor({ path, blockSize, maxBlockCount }) {
        this.logger = new Logger('subframe/file-storage')
        this.logger.debug('Initializing FileStorage...')
        this.path = path
        this.filehandle = null
        this.blockSize = blockSize
        this.occupiedBlocks = [0]
        this.maxBlockCount = maxBlockCount
        this._loadOrCreate()
    }

    async _loadOrCreate() {
        if (!fs.existsSync(this.path)) {
            this.logger.debug(`File ${this.path} does not exist. Creating...`)
            await this._create()
        } else {
            try {
                this.logger.debug(`File ${this.path} exits. Loading...`)
                await this._load()
            } catch (e) {
                this.logger.fatal(`Failed to load file at '${this.path}': ${e.toString()}`)
            }
        }
    }

    async _create() {
        this.logger.debug('Creating new data file...')
        fs.mkdirSync(path.dirname(this.path), { recursive: true })
        fs.writeFileSync(this.path, null)

        this.filehandle = await fs.promises.open(this.path, 'r+')

        await this._writeBlockZero()
        this.logger.debug('Done.')
    }

    async _load() {
        this.filehandle = await fs.promises.open(this.path, 'r+')
        let buf = Buffer.alloc(16)
        await this.filehandle.read(buf, 0, buf.length, 0)
        const blockZeroLength = buf.toString().split(';')[0]
        buf = Buffer.alloc(parseInt(blockZeroLength))
        await this.filehandle.read(buf, 0, buf.length, blockZeroLength.length + 1)
        const blockZero = JSON.parse(buf.toString())
        this.occupiedBlocks = blockZero.occupiedBlocks
        this.blockSize = blockZero.blockSize
    }

    async _writeBlockZero() {
        const { occupiedBlocks, blockSize } = this
        const data = JSON.stringify({
            occupiedBlocks: occupiedBlocks.sort(),
            blockSize
        })
        const dataLength = data.length

        this._writeBlockToFile(0, `${dataLength};${data}`)
    }

    async _writeBlock(offset, data) {
        await this._writeBlockToFile(offset, data)
        this.occupiedBlocks.filter(occupiedBlock => occupiedBlock !== offset).push(offset)
        await this._writeBlockZero()
    }

    async _writeBlockToFile(offset, data) {
        const buf = Buffer.from(data)
        if (buf.length > this.blockSize)
            throw new NamedError(errorTypes.fileStorage.dataExceedsBlockSize, `Data does not fit into block size (Length: ${buf.length}, block size: ${this.blockSize})`)

        await this.filehandle.write(buf, 0, buf.length, offset * this.blockSize)
    }

    async _readBlock(offset) {

    }

    async _freeBlock(offset) {
        const buf = Buffer.alloc(this.blockSize)
        await this.filehandle.write(buf, 0, buf.length, offset * this.blockSize)
        this.occupiedBlocks = this.occupiedBlocks.filter(occupiedBlock => occupiedBlock !== offset)
        return this._writeBlockZero()
    }

    _getHighestBlockOffset() {
        return this.occupiedBlocks.sort().pop();
    }
}

module.exports = FileStorage