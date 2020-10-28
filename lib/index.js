const fs = require("fs");
const path = require("path");
const getDir = require("./parse");

/**
 * 务必配置src的alias为@
 * 自定义路由映射表生成
 */
class UxAutoRouterPlugin {
  constructor({ pagePath, output, filename }) {
    if (!pagePath) {
      throw new Error("`pagePath not undefined`");
    }
    this.filename = filename || "router.js";
    //this.srcAlias = srcAlias || "src";
    this.pagePath = pagePath;
    this.output = output || ".";
    this.componentBase = path.basename(pagePath);
    this.route = null;
  }

  apply(compiler) {
    // watchRun 是异步 hook，使用 tapAsync 触及它，还可以使用 tapPromise/tap(同步)
    // 在自动编译前生成路由映射表
    compiler.hooks.watchRun.tapAsync(
      "RouterPlugin",
      (compilation, callback) => {
        this.buildRouter();
        callback();
      }
    );

    // 在构建项目时生成路由映射表
    compiler.hooks.beforeRun.tapAsync(
      "RouterPlugin",
      (compilation, callback) => {
        this.buildRouter();
        callback();
      }
    );
  }

  buildRouter() {
    // import的静态路由
    const staticRoute = [];
    // 读取文件目录生成路由对象
    const routers = getDir(this.output, this.pagePath, staticRoute);
    // 检测是否需要重新生成路由，即pages是否发生了改变，此处做的比较简单只是转换成JSON字符串进行比对
    if (!this.route || JSON.stringify(this.route) !== JSON.stringify(routers)) {
      // 格式化字符串
      let str = JSON.stringify([routers], null, 4);
      // JSON文件转换换成js文件
      str = str.replace(/"/g, "");

      console.log("[自动写入路由配置，成功]");
      // 更新路由JSON字符串
      this.route = routers;
      staticRoute.unshift("import loadable from '@loadable/component';\n");
      fs.writeFileSync(
        `${this.output}/${this.filename}`,
        `${staticRoute.map((item) => item).join("")}\n\n//路由映射表\nconst router = ${str};\n\nexport default router;`,
        { flag: "w", encoding: "utf-8", mode: "0666" }
      );
    }
  }
}

module.exports = UxAutoRouterPlugin;
