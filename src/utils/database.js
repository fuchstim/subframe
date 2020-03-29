const fs = require('fs')
const path = require('path')
const AsyncLock = require('async-lock')

const { Logger } = require('../logger')
const { NamedError, errorTypes } = require('../errors')

class Database {
    constructor({ path }) {
        this.logger = new Logger('subframe/database')
        this.logger.info(`'Initializing database at '${path}'..`)
        this.path = path
        this.hashmap = null
        this.lock = new AsyncLock()
        this._loadOrCreate()
    }

    _loadOrCreate() {
        if (!fs.existsSync(this.path)) {
            this.logger.debug(`File ${this.path} does not exist. Creating...`)
            this._create()
        } else {
            try {
                this.logger.debug(`File ${this.path} exits. Loading...`)
                this._load()
            } catch (e) {
                this.logger.fatal(`Failed to load database file at '${this.path}': ${e.toString()}`)
            }
        }
    }

    _create() {
        this.logger.debug('Creating new database...')
        this.hashmap = {}
        fs.mkdirSync(path.dirname(this.path), { recursive: true })
        this._save()
        this.logger.debug('Done.')
    }

    _load() {
        this.logger.debug('Loading database from file...')
        try {
            this.hashmap = JSON.parse(fs.readFileSync(this.path, { encoding: 'utf-8' }))
            this.logger.debug(`Done. Loaded ${Object.keys(this.hashmap).length} entries.`)
        } catch (e) {
            this.logger.fatal(`Failed to load database from file: ${e.toString()}`)
        }
    }

    _save() {
        this.logger.debug('Awaiting file lock...')
        this.lock.acquire('file', () => {
            this.logger.debug('Writing database to file...')
            fs.writeFileSync(this.path, JSON.stringify(this.hashmap), { encoding: 'utf-8' })
        }, () => {
            this.logger.debug('Wrote database to file')
        })
    }

    write(key, data) {
        this.hashmap[key] = data
        this._save()
    }

    read(key) {
        if (!this.hasRecord(key))
            throw new NamedError(errorTypes.database.missingRecord, `Unknown key: ${key}`)

        return this.hashmap[key]
    }

    hasRecord(key) {
        return this.hashmap.hasOwnProperty(key)
    }
}

module.exports = Database