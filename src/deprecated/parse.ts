import { titleCase, isRoute, replaceDiagonal } from "./name";
import * as fs from "fs";
import * as path from "path";

//递归遍历文件夹
/**
 *
 * @param {string} output 输出文件夹
 * @param {string} dir 当前文件夹路径
 * @param {string} staticRoute 静态路由列表
 * @param {any} routePath 路由延伸路径
 */
function getDir(
    output: string,
    dir: string,
    staticRoute: string[],
    routePath: any
) {
    routePath = routePath || "";
    try {
        const files = fs.readdirSync(dir);
        // 遍历子集路由
        const childs: string[] = [];
        // 统一路由配置
        let routeConfig: UxPlugin.RouteConfig = {
            component: "",
            path: "",
            child: [],
        };
        files.forEach(async (filename: string) => {
            let config;
            // 是否有配置文件
            let hasConfig = false;
            let filePath = path.join(dir, filename);
            const stat = fs.statSync(filePath);
            // 解析路由配置文件

            if ("route" === path.basename(filePath, ".config")) {
                config = JSON.parse(fs.readFileSync(filePath) + "");
                hasConfig = true;

                //如果配置了非动态路由
                if (config.hasOwnProperty("noLazy") && config.noLazy) {
                    // 如果没有component对象就对其赋值静态路由
                    routeConfig.component =
                        "Page" + titleCase(path.basename(routePath));
                    // 在头部导入静态路由
                    staticRoute.push(
                        `import Page${titleCase(
                            path.basename(routePath)
                        )} from '${replaceDiagonal(
                            path.relative(output, dir)
                        )}/index';\n`
                    );
                }
            }

            // 解析路由根组件
            if (
                "index" === path.basename(filePath, ".tsx") ||
                "index" === path.basename(filePath, ".jsx")
            ) {
                // 默认设置为动态路由，此处必须使用loadable与import这样用
                // 因为如果在组件中传递参数来指引路由组件会无效，因为import内部是使用一组正则来判断，无法使用变量
                // 只能生成既定的值来作为参数，loadable也需要加在这，因为如果不这样做也会导致组件不能正常加载
                if (!routeConfig.component) {
                    routeConfig.component = `loadable(()=>import('${replaceDiagonal(
                        path.relative(output, dir)
                    )}/${path.basename(filePath).split(".")[0]}'))`;
                }
                routeConfig.path = `'${routePath}'`;
            }
            if (stat.isDirectory() && isRoute(filename)) {
                childs.push(filePath);
            }

            //如果有配置文件route.config 则直接添加到路由配置中
            if (hasConfig) {
                Object.assign(config, routeConfig);
                routeConfig = config;
            }
        });

        childs.forEach((filePath) => {
            let child = getDir(
                output,
                filePath,
                staticRoute,
                routePath + "/" + path.basename(filePath)
            );
            routeConfig.child.push(child);
        });

        return routeConfig;
    } catch (e) {
        console.log("×[自动写入路由配置，失败]");
        throw new Error(e);
    }
}

module.exports = getDir;
