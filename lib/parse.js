const fs = require("fs");
const path = require("path");
const { titleCase, isRoute, replaceDiagonal } = require("./utils");

//递归遍历文件夹
/**
 *
 * @param {string} output 输出文件夹
 * @param {string} dir 当前文件夹路径
 * @param {string} staticRoute 静态路由列表
 * @param {any} routePath 路由延伸路径
 */
function getDir(output, dir, staticRoute, routePath) {
  routePath = routePath || "";
  try {
    const files = fs.readdirSync(dir);
    // 遍历子集路由
    const childs = [];
    // 统一路由配置
    let routeConfig = {
      child: [],
    };
    files.forEach(async (filename) => {
      let config;
      // 是否有配置文件
      let hasConfig = false;
      let pathname = path.join(dir, filename);
      const stat = fs.statSync(pathname);
      // 解析路由配置文件
      if ("route" === path.basename(pathname, ".config")) {
        config = JSON.parse(fs.readFileSync(pathname));
        hasConfig = true;

        //如果配置了非动态路由
        if (config.hasOwnProperty("noLazy") && config.noLazy) {
          // 如果没有component对象就对其赋值静态路由
          routeConfig.component = "Page" + titleCase(path.basename(routePath));
          // 在头部导入静态路由
          staticRoute.push(
            `import Page${titleCase(
              path.basename(routePath)
            )} from '${replaceDiagonal(path.relative(output, dir))}/index';\n`
          );
        }
      }

      // 解析路由根组件
      if (
        "index" === path.basename(pathname, ".tsx") ||
        "index" === path.basename(pathname, ".jsx")
      ) {
        // 默认设置为动态路由，此处必须使用loadable与import这样用
        // 因为如果在组件中传递参数来指引路由组件会无效，因为import内部是使用一组正则来判断，无法使用变量
        // 只能生成既定的值来作为参数，loadable也需要加在这，因为如果不这样做也会导致组件不能正常加载
        if (!routeConfig.component) {
          routeConfig.component = `loadable(function (){return import('${replaceDiagonal(
            path.relative(output, dir)
          )}/${path.basename(pathname)}')})`;
        }
        routeConfig.path = `'${routePath}'`;
      } else if (stat.isDirectory() && isRoute(filename)) {
        childs.push(pathname);
      }

      //如果有配置文件route.config 则直接添加到路由配置中
      if (hasConfig) {
        Object.assign(config, routeConfig);
        routeConfig = config;
      }
    });

    childs.forEach((pathname) => {
      let child = getDir(
        output,
        pathname,
        staticRoute,
        routePath + "/" + path.basename(pathname)
      );
      routeConfig.child.push(child);
    });

    return routeConfig;
  } catch (e) {
    console.log("×[自动写入路由配置，失败]");
    throw new Error(e)
  }
}

module.exports = getDir;
