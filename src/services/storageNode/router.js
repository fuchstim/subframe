const express = require('express')
const router = express.Router()

const StorageHandler = require('./handler')

router.get('/get-info/:messageID', (req, res, next) => StorageHandler.getFileInfo({ req, res, next }))

module.exports = router