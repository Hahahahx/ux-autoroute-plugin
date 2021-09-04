declare module "ux-autoroute-plugin" {
  export class AutoRouterWebPackPlugin {
    constructor(params: PluginParams);
  }

  export function AutoRouterVitePlugin(params: PluginParams): any;

  interface RouteConfig {
    component: string;
    path: string;
    child: RouteConfig[];
  }

  interface PluginParams {
    // 输入路径
    pagePath: string;
    // 输出路径
    output: string;
    // 输出文件名
    filename: string;
    // 默认关闭按需加载import
    defaultLazyImport?: boolean;
  }
}
