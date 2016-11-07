var path = require('path')

module.exports = {
	entry: './static/index.js',
	output: {
		path: path.resolve(__dirname, 'bundle'),
		filename: './bundle.js'
	},
	module: {
		loaders: [
			{
				test: /\.js$/,
				loader: 'babel-loader',
				query:
				{
					presets:['es2015']
				}
			}
		]
	}
}
