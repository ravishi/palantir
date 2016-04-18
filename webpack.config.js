var path = require('path');
var webpack = require('webpack');

module.exports = {
    entry: {
        index: "./client/index.html",
        bundle: [
            // enable hot-reloading for IE
            'eventsource-polyfill',
            'react-hot-loader/patch',
            './client/index.jsx'
        ]
    },
    output: {
        path: path.join(__dirname, 'build'),
        filename: '[name].js'
    },
    plugins: [
        new webpack.HotModuleReplacementPlugin(),
        new webpack.NoErrorsPlugin()
    ],
    module: {
        loaders: [
            {test: /\.jsx?$/, exclude: /node_modules/, loader: 'babel'},
            {test: /\.html?$/, loader: 'file?name=[name].[ext]'}
        ]
    }
};
