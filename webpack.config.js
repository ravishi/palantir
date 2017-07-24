const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: [
        'eventsource-polyfill',
        './client/index.jsx'
    ],
    output: {
        path: path.join(__dirname, 'build'),
        filename: '[name].js'
    },
    plugins: [
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NoEmitOnErrorsPlugin(),
        new HtmlWebpackPlugin({
            template: 'client/index.ejs'
        })
    ],
    module: {
        loaders: [
            {test: /\.jsx?$/, exclude: /node_modules/, loader: 'babel-loader'},
            {test: /\.jpg$/, loader: 'file-loader?name=[name].[ext]'}
        ]
    }
};
