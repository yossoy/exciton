const UglifyJsPlugin = require('uglifyjs-webpack-plugin')

module.exports = {
  mode: 'production',
  devtool: 'source-map',
  entry: [
    './js/sample.js',
  ],
  optimization: {
    usedExports: true,
    concatenateModules: true,
    occurrenceOrder: true
  },
  output: {
    path: `${__dirname}/resources`,
    filename: 'sample.js',
    library: 'sample',
    libraryTarget: 'commonjs2'
  },
  module: {rules: [{test: /\.js$/, exclude: /node_modules/, use: {loader: 'babel-loader'}}]},
  optimization: {
    minimizer: [new UglifyJsPlugin({
      uglifyOptions: {
        safari10: true,
        compress: {
          drop_console: false,
        }
      }
    })],
  }
};