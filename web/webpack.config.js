const path = require("path");
const ExtractTextPlugin = require("extract-text-webpack-plugin");


const extractSASS = new ExtractTextPlugin({ filename: 'bundle.css' })
const SRC_DIR = path.resolve(__dirname, 'src');
const BUILD_DIR = path.resolve(__dirname, 'static/build');

module.exports = {
    entry: './src/index.jsx',
    output: {
		path: BUILD_DIR,
		filename: 'bundle.js',
		publicPath: "/",
	},
	module: {
		rules: [
		    {
				test: /\.jsx$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader"
				}
			},
			{
				test: /\.scss$/,
				use: extractSASS.extract({
					fallback: 'style-loader',
					use: [ 'css-loader', 'sass-loader' ]
				})
			}
		]
	},
	resolve: {
		extensions: ['.js', '.jsx'],
	},
	plugins: [
		extractSASS
	]
};