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
    }
}