# ux-autoroute-plugin

webpack 自动路由插件，该插件用于生成路由映射表文件

```Javascript
    // paths.appSrc 为src目录
    new UxAutoRouterPlugin({
        pagePath: path.join(paths.appSrc, 'pages'),
        output: path.join(paths.appSrc, 'config'),
        srcAlias:'@',            // 生成的文件内的src别名，需要与webpack配合，如果不加这个配置默认为src
        filename:'routerConfig.ts'
    })

```

在 src 下读取 pages 文件夹，在 config 文件夹下，生成文件 routerConfig.ts：

```Javascript
import Page from 'src/pages/index'; //配置srcAlias:'@'，就为'@/pages/index'
export const routeConfig = [
    {
        noLazy: true,
        child: [
            {
                default: true,
                child: [],
                componentPath: 'pages/login/index.tsx',
                path: '/login'
            },
            {
                child: [],
                componentPath: 'pages/main/index.tsx',
                path: '/main'
            }
        ],
        componentPath: 'pages/index.tsx',
        path: '',
        component: Page
    }
]
```

路由要求所有路由文件夹都是 `小写且不含特殊字符` ，且文件夹内必须有 index.jsx|index.tsx 为该路由的路由组件。

noLazy，路由默认都会生成为动态路由，会进行代码分片，加上这个配置之后的路由不会生成动态路由，而是使用 import 静态引入

default，默认路由，项目会以 pages 为根路由即`/`，往下都是其子路由，如果不使用根路由，直接将其视作引导页，其子路由中需要在某一个路由配置中加入 default，<br/>
如在 pages 文件夹下新建 login 文件夹，即 login 路由，添加配置 route.config，那么在访问`/`会被引导到`/login`。

pages/login/route.config，`/login`路由的配置文件 route.config

```Json
{
    "default":true
}
```

如此在 pages，根路由下一般会取消动态路由，因为作为跟路由代码量本身没有必要分片加载，所以在配置中加入 noLazy：
pages/route.config，`/`路由配置文件：

```Json
{
    "noLazy":true
}
```

--------------------------------------

本插件始终只作为生成路由映射表文件，对映射表的解析需要自行定义，如重新封装react-router