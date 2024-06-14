const path = require('path')
const qiniu = require('qiniu')
const walkSync = require('walk-sync')
const chalk = require('chalk')

// target directory
const targetDir = './api'

// access key
const accessKey = process.env.QINIU_ACCESS_KEY

// secret key
const secretKey = process.env.QINIU_SECRET_KEY

// bucket
const bucket = process.env.QINIU_BUCKET

// config
const config = new qiniu.conf.Config()

// zone
config.zone = qiniu.zone[process.env.QINIU_ZONE]

function uploadFile(localFile, key) {
	// options
	const options = {
		scope: `${bucket}:${key}`,
	}

	// mac
	const mac = new qiniu.auth.digest.Mac(accessKey, secretKey)

	// put policy
	const putPolicy = new qiniu.rs.PutPolicy(options)

	// upload token
	const uploadToken = putPolicy.uploadToken(mac)

	return new Promise((resolve, reject) => {
		const formUploader = new qiniu.form_up.FormUploader(config)
		const putExtra = new qiniu.form_up.PutExtra()
		formUploader.putFile(uploadToken, key, localFile, putExtra, function (respErr, respBody, respInfo) {
			if (respErr) {
				throw respErr
			}
			if (respInfo.statusCode === 200) {
				console.log(`${chalk.green('uploaded')} ${localFile} => ${key}`)
				resolve()
			} else if (respInfo.statusCode === 614) {
				console.log(`${chalk.yellow('exists')} ${localFile} => ${key}`)
				resolve()
			} else {
				const errMsg = `${chalk.red('error[' + respInfo.statusCode + ']')} ${localFile} => ${key}`
				console.error(errMsg)
				reject(new Error(respBody))
			}
		})
	})
}

async function main() {
	// paths
	const paths = walkSync(targetDir, {
		includeBasePath: true, directories: false,
	})

	// iterate paths
	for (const filePath of paths) {
		const localFile = path.resolve(filePath)
		const key = filePath.replace(targetDir + '/', '')
		try {
			await uploadFile(localFile, key)
		} finally {
			// do nothing
		}
	}
}

(async () => {
	await main()
})()
