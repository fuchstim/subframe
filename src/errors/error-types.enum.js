module.exports = {
    unknown: {
        code: 'error/unknown',
        httpCode: 500,
        description: 'An unknown error occured.'
    },

    request: {
        badRequest: {
            code: 'request/bad-request',
            httpCode: 400,
            description: 'Bad Request'
        },
        unknownResource: {
            code: 'request/unknown-resource',
            httpCode: 404,
            description: 'Unknown Resource'
        },
    },

    database: {
        missingRecord: {
            code: 'database/missing-record',
            httpCode: 500,
            description: 'Failed to retrieve record for unknown key'
        },

        writeError: {
            code: 'database/write-error',
            httpCode: 500,
            description: 'Failed to write database file to disk'
        },

        readError: {
            code: 'database/read-error',
            httpCode: 500,
            description: 'Failed to read database file from disk'
        }
    },

    fileStorage: {
        dataExceedsBlockSize: {
            code: 'file-storage/exceeds-block-size',
            httpCode: 400,
            description: 'Data length exceeds block size'
        }
    }
}