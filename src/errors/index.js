const errorTypes = require('./error-types.enum')

class NamedError extends Error {
    constructor(errorType = errorTypes.unknown, message, additionalInformation = {}) {
        super(message)

        this.code = errorType.code
        this.httpCode = errorType.httpCode
        this.description = additionalInformation.description || errorType.description
        this.name = `NamedError [${this.code}]`
        this.message = message
        this.additionalInformation = additionalInformation

        Error.captureStackTrace(this, this.constructor)
    }

    toString = (prependErrorCode = true) => (prependErrorCode ? `[${this.code}] ` : '') + `${this.message}`
}

module.exports = {
    NamedError,
    errorTypes
}