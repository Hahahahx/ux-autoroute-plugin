declare module "ux-autoroute-plugin" {
  export class AutoRouterWebPackPlugin {
    constructor(params: PluginParams);
  }

  export function AutoRouterVitePlugin(params: PluginParams);

  interface RouteConfig {
    component: string;
    path: string;
    child: RouteConfig[];
  }

  interface PluginParams {
    pagePath: string;
    output: string;
    filename: string;
  }
}
