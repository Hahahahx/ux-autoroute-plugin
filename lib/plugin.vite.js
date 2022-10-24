"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AutoRouterVitePlugin = void 0;
const path = require("path");
const process = require("child_process");
function AutoRouterVitePlugin(params) {
    return {
        name: "RouterPlugin",
        buildStart() {
            buildRouter(params);
        },
        handleHotUpdate() {
            buildRouter(params);
        },
    };
}
exports.AutoRouterVitePlugin = AutoRouterVitePlugin;
function buildRouter(params) {
    const parse_exe = path.join(__dirname, "..", "parse.exe");
    process.exec(`${parse_exe} parse -o ${params.output} -r ${params.pagePath} -n ${params.filename} -i "${params.defaultLazyImport}"`, (err, stdout, stderr) => {
        // console.log(stdout);
    });
}
// module.exports = { AutoRouterVitePlugin };
