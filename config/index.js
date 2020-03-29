const environment = process.env.NODE_ENV

const defaults = require('./defaults.json')
const development = require('./development.json')

const config = {
    development
};

module.exports = {
    ...defaults,
    ...config[environment]
};