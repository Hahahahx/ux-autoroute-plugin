"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AutoRouterWebPackPlugin = void 0;
const path = require("path");
const process = require("child_process");
/**
 * 务必配置src的alias为@
 * 自定义路由映射表生成
 */
class AutoRouterWebPackPlugin {
    constructor({ pagePath, output, filename }) {
        if (!pagePath) {
            throw new Error("`pagePath not undefined`");
        }
        this.filename = filename || "router.js";
        this.pagePath = pagePath;
        this.output = output || ".";
    }
    apply(compiler) {
        // watchRun 是异步 hook，使用 tapAsync 触及它，还可以使用 tapPromise/tap(同步)
        // 在自动编译前生成路由映射表
        compiler.hooks.watchRun.tapAsync("RouterPlugin", (compilation, callback) => {
            this.buildRouter();
            callback();
        });
        // 在构建项目时生成路由映射表
        compiler.hooks.beforeRun.tapAsync("RouterPlugin", (compilation, callback) => {
            this.buildRouter();
            callback();
        });
    }
    buildRouter() {
        const parse_exe = path.join(__dirname, "parse.exe");
        process.exec(`${parse_exe} parse -o ${this.output} -r ${this.pagePath} -n ${this.filename}`, (err, stdout, stderr) => {
            // console.log(stdout);
        });
        // const { status } = process.spawnSync(
        //     parse_exe,
        //     [
        //         "parse",
        //         "-o " + this.output,
        //         "-r " + this.pagePath,
        //         "-n " + this.filename,
        //     ],
        //     { stdio: "inherit" }
        // );
        // (process as any).exitCode = status === null ? 1 : status;
    }
}
exports.AutoRouterWebPackPlugin = AutoRouterWebPackPlugin;
