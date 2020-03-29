const { NamedError, errorTypes: { request: { unknownResource } } } = require('../errors')
module.exports = (req, res) => {
    if (!res.data && !res.error)
        res.error = new NamedError(unknownResource, 'Unknown Endpoint')

    if (res.error) {
        res.statusCode = res.error.httpCode
        res.data = res.error.description
        res.logger.error(res.error)
    }

    res.status(res.statusCode || 200)
    res.json({
        status: res.statusCode || 200,
        data: res.data
    })
    res.end()
}