const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');

const app = {
  mode: 'production',
  entry: ['./src/exciton.js'],
  optimization:
      {usedExports: true, concatenateModules: true, occurrenceOrder: true},
  output: {path: `${__dirname}/data`, filename: 'exciton.js'},
  module: {
    rules: [
      {test: /\.js$/, exclude: /node_modules/, use: {loader: 'babel-loader'}},
      {test: /\.css$/, use: [MiniCssExtractPlugin.loader, 'css-loader']}
    ]
  },
  optimization: {
    minimizer: [
      new UglifyJsPlugin({
        uglifyOptions: {
          safari10: true,
          compress: {
            drop_console: true,
          },
          sourceMap: false
        },
      }),
      new OptimizeCSSAssetsPlugin({})
    ],
  },
  plugins: [new MiniCssExtractPlugin({filename: '../data/default.css'})]
};


module.exports = [app];
