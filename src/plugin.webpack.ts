import * as path from "path";
import * as process from "child_process";

/**
 * 务必配置src的alias为@
 * 自定义路由映射表生成
 */
export class AutoRouterWebPackPlugin {
    public output: string;
    public pagePath: string;
    public filename: string;

    constructor({ pagePath, output, filename }: UxPlugin.PluginParams) {
        if (!pagePath) {
            throw new Error("`pagePath not undefined`");
        }
        this.filename = filename || "router.js";
        this.pagePath = pagePath;
        this.output = output || ".";
    }

    apply(compiler: any) {
        // watchRun 是异步 hook，使用 tapAsync 触及它，还可以使用 tapPromise/tap(同步)
        // 在自动编译前生成路由映射表
        compiler.hooks.watchRun.tapAsync(
            "RouterPlugin",
            (compilation: any, callback: any) => {
                this.buildRouter();
                callback();
            }
        );

        // 在构建项目时生成路由映射表
        compiler.hooks.beforeRun.tapAsync(
            "RouterPlugin",
            (compilation: any, callback: any) => {
                this.buildRouter();
                callback();
            }
        );
    }

    buildRouter() {
        const parse_exe = path.join(__dirname,  "parse.exe");
        process.exec(
            `${parse_exe} parse -o ${this.output} -r ${this.pagePath} -n ${this.filename}`,
            (err:any, stdout:any, stderr:any) => {
                // console.log(stdout);
            }
        );
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
