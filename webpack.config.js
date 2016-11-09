var path = require('path');

module.exports = {
	devtool: 'source-map',
	entry: './static/index.js',
	output: {
		path: path.resolve(__dirname, 'bundle/assets'),
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
			},
			{ test: /\.less$/, loader: 'style!css!less' },
			{ test: /\.(png|woff|woff2|eot|ttf|svg)$/, loader: 'url-loader?limit=100000' }
		]
	},
	devServer: {
		publicPath: '/assets',
		proxy: {
			'/api': {
				target: 'http://localhost:6288',
				secure: false
			}
		}
	}
};
