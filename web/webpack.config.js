const path = require("path");

const TerserJSPlugin = require('terser-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const BUILD_DIR = path.resolve(__dirname, './static/build');

module.exports = {
	entry: './src/index.jsx',
 	optimization: {
		minimizer: [new TerserJSPlugin({})]
	},
	output: {
		path: BUILD_DIR,
		filename: '[name].[contenthash].js',
		publicPath: '/static/build',
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
				use: [MiniCssExtractPlugin.loader, 'css-loader', 'sass-loader']
			}
		]
	},
	resolve: {
		extensions: ['.js', '.jsx'],
	},
	plugins: [
 		new CleanWebpackPlugin(),
		new MiniCssExtractPlugin({
			filename: '[name].[contenthash].css'
		}),
		new HtmlWebpackPlugin({
			filename: path.join(__dirname, './views/index.html'),
			template: path.join(__dirname, './views/base.html'),
		}),
	]
};
