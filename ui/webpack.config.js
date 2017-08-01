'use strict';

const path = require('path');
const autoprefixer = require('autoprefixer');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = (env) => {
    const src = path.join(__dirname, 'src');
    const dist = path.join(__dirname, 'dist');

    return {
        entry: ['core-js/shim', 'whatwg-fetch', './src/index.tsx'],

        output: {
            //filename: 'bundle.[chunkhash].js',
            filename: 'bundle.js',
            path: dist,
        },

        resolve: {
            modules: [
                'node_modules'
            ],
            extensions: [
                '.js',
                '.ts',
                '.tsx',
                '.scss',
                '.pug',
                '.json'
            ]
        },

        module: {
            rules: [{
                test: /\.(tsx|ts)?$/,

                exclude: /node_modules/,
                include: [
                    src,
                ],
                use: [
                    'react-hot-loader',
                    {
                        loader: 'ts-loader',
                        options: {
                            compilerOptions: {
                                target: env === 'prod' ? 'es5' : 'es2017'
                            }
                        }
                    }
                ]
            }, {
                test: /\.scss$/,
                use: [
                    'style-loader',
                    'css-loader',
                    {
                        loader: 'postcss-loader',
                        options: {
                            plugins: function () {
                                return [autoprefixer({
                                    browsers: ["last 2 versions", "ie 9"]
                                })];
                            }
                        }
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            includePaths: [path.join(src, 'style')],
                            outputStyle: 'expanded'
                        }
                    }
                ]
            },]
        },

        devtool: 'eval-source-map',

        plugins: [
            new HtmlWebpackPlugin({
                filename: 'index.html',
                template: 'template/index.html',
                chunks: ['bundle']
            }),
            new webpack.DefinePlugin({
                ENDPOINT: env === 'prod' ? '\"/data.json\"' : null
            })
        ]
    }
};
