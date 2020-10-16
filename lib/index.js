import "core-js/modules/es.array.join";
import "core-js/modules/es.array.map";
import "core-js/modules/es.regexp.exec";
import "core-js/modules/es.string.replace";
import _classCallCheck from "@babel/runtime/helpers/esm/classCallCheck";
import _createClass from "@babel/runtime/helpers/esm/createClass";

var fs = require("fs");

var path = require("path");

var getDir = require("./utils");
/**
 * 务必配置src的alias为@
 * 自定义路由映射表生成
 */


var UxAutoRouterPlugin = /*#__PURE__*/function () {
  function UxAutoRouterPlugin(_ref) {
    var pagePath = _ref.pagePath,
        output = _ref.output,
        srcAlias = _ref.srcAlias,
        filename = _ref.filename;

    _classCallCheck(this, UxAutoRouterPlugin);

    if (!pagePath) {
      throw new Error("`pagePath not undefined`");
    }

    this.filename = filename || "router.js";
    this.srcAlias = srcAlias || "src";
    this.pagePath = pagePath;
    this.output = output || ".";
    this.componentBase = path.basename(pagePath);
    this.route = null;
  }

  _createClass(UxAutoRouterPlugin, [{
    key: "apply",
    value: function apply(compiler) {
      var _this = this;

      // watchRun 是异步 hook，使用 tapAsync 触及它，还可以使用 tapPromise/tap(同步)
      // 在自动编译前生成路由映射表
      compiler.hooks.watchRun.tapAsync("RouterPlugin", function (compilation, callback) {
        _this.buildRouter();

        callback();
      }); // 在构建项目时生成路由映射表

      compiler.hooks.beforeRun.tapAsync("RouterPlugin", function (compilation, callback) {
        _this.buildRouter();

        callback();
      });
    }
  }, {
    key: "buildRouter",
    value: function buildRouter() {
      // import的静态路由
      var staticRoute = []; // 读取文件目录生成路由对象

      var routers = getDir(this.srcAlias, this.pagePath, this.componentBase, staticRoute); // 检测是否需要重新生成路由，即pages是否发生了改变，此处做的比较简单只是转换成JSON字符串进行比对

      if (!this.route || JSON.stringify(this.route) !== JSON.stringify(routers)) {
        // 格式化字符串
        var str = JSON.stringify([routers], null, 4); // JSON文件转换换成js文件

        str = str.replace(/"/g, "");
        console.log("[自动写入路由配置，成功]"); // 更新路由JSON字符串

        this.route = routers;
        fs.writeFileSync(this.output + "/" + this.filename, staticRoute.map(function (item) {
          return item;
        }).join("") + "export const routeConfig = " + str, {
          flag: "w",
          encoding: "utf-8",
          mode: "0666"
        });
      }
    }
  }]);

  return UxAutoRouterPlugin;
}();

module.exports = UxAutoRouterPlugin;