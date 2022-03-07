package project

import (
	"MockConfig/module"
	"MockConfig/tool"
	"bytes"
	"io/ioutil"
	"os"
)

var gitignore = `*/.DS_Store
.idea
node_modules
dist
android
build
webpack/node_modules
webpack/package-lock.json
`

func es6Project(projectName string, path string) {

	writeGitignore(path, gitignore)

	writeAppConfig(projectName, projectName,module.ES6)

	mkSRCdir(path)

	mkWebpackDir(path)

}

func mkWebpackDir(path string) {
	if err := os.Mkdir(path+"/webpack", os.ModePerm); err == nil {
	} else {
		panic(err)
	}

	if err := ioutil.WriteFile(path+"/webpack/package.json", bytes.NewBufferString(packageJson).Bytes(), os.ModePerm); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(path+"/webpack/webpack.config.js", bytes.NewBufferString(webpackConfig).Bytes(), os.ModePerm); err != nil {
		panic(err)
	}

	tool.ExecCmd("npm","i","--prefix",path+"/webpack")

}

var webpackConfig = `const path = require('path');

const fs = require('fs-extra');

let json = fs.readJsonSync(path.resolve(__dirname, '../tiny.json'))

let pages = json.runtime.pages;

let entry = {}

let projectPath = path.resolve(__dirname, "../")

for (let i = 0; i < pages.length; i++) {
    entry[pages[i].name] = getAbsPath(projectPath, pages[i].source);
}

module.exports = {
    mode: "development",
    devtool: 'inline-source-map',
    entry: entry,
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, '../build'),
        clean: false,
    },
    resolve: {
        extensions: ['.js'],
    },
    optimization: {
        usedExports: false
    }
}


function getAbsPath(project, source) {
    if ("./" === source.substring(0, 2)) {
        return project + "/" + source.substring(2)
    } else if ("../" === source.substring(0, 3)) {
        let count = 0
        let lastIndex = 0
        for (let i = 0; i < source.length; i++) {
            if (source[i] === ".") {
                if (i + 1 < source.length && source[i + 1] == ".") {
                    if (i + 2 < source.length && source[i + 2] == "/") {
                        count++;
                        lastIndex = i + 2;
                        i = i + 2;
                    }
                }
            }
        }
        for (let c = count; c > 0; c--) {
            project = project.substring(0, mLastIndex(project))
        }
        return project + source.substring(lastIndex)
    } else if ("/" === source.substring(0, 1)) {
        return source
    }
}


function mLastIndex(str) {
    for (let c = str.length; c >= 0; c--) {
        if ("/" === str[c]) {
            return c;
        }
    }
}
`

var packageJson = `{
  "scripts": {
    "build": "webpack --mode production",
    "watch": "webpack --watch"
  },
  "license": "ISC",
  "dependencies": {
    "webpack": "^5.58.2",
    "webpack-cli": "^4.9.0"
  },
  "devDependencies": {
    "fs-extra": "^10.0.0",
    "webpack-dev-server": "^4.7.4"
  }
}`
