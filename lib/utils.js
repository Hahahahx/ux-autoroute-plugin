import "core-js/modules/es.array.concat";
import "core-js/modules/es.array.for-each";
import "core-js/modules/es.array.join";
import "core-js/modules/es.array.slice";
import "core-js/modules/es.object.assign";
import "core-js/modules/es.regexp.exec";
import "core-js/modules/es.string.split";
import "core-js/modules/web.dom-collections.for-each";

//递归遍历文件夹
function getDir(srcAlias, dir, componentBase, staticRoute, routePath) {
  routePath = routePath || ""; // console.log('[开始]   ======>   ' + dir)

  try {
    var files = fs.readdirSync(dir);
    var childs = [];
    var router = {
      child: []
    };
    files.forEach(function (filename) {
      var component;
      var hasComponent = false;
      var pathname = path.join(dir, filename); // console.log('[读取]   ======>   ' + pathname)

      var stat = fs.statSync(pathname); // 解析路由配置文件

      if ("route" === path.basename(pathname, ".config")) {
        component = JSON.parse(fs.readFileSync(pathname));
        hasComponent = true; //console.log('==================>    ', component, Object.hasOwnProperty(component, 'noLazy'))
        // 判断如果非动态引入组件，就会直接解析成import，静态引入组件，并添加到component属性中

        if (component.hasOwnProperty("noLazy") && component.noLazy) {
          router.component = "Page" + titleCase(path.basename(routePath));
          staticRoute.push("import Page".concat(titleCase(path.basename(routePath)), " from '").concat(srcAlias, "/").concat(componentBase + routePath, "/index';\n"));
        }
      } // 解析路由根组件


      if ("index" === path.basename(pathname, ".tsx") || "index" === path.basename(pathname, ".jsx")) {
        router.componentPath = "'" + componentBase + routePath + "/" + path.basename(pathname) + "'";
        router.path = "'" + routePath + "'";
      } else if (stat.isDirectory() && isRoute(filename)) {
        childs.push(pathname);
      }

      if (hasComponent) {
        Object.assign(component, router);
        router = component;
      }
    });
    childs.forEach(function (pathname) {
      var child = getDir(srcAlias, pathname, componentBase, staticRoute, routePath + "/" + path.basename(pathname));
      router.child.push(child);
    });
    return router;
  } catch (e) {
    console.log("×[自动写入路由配置，失败]");
    return {};
  }
}

function titleCase(str) {
  if (str[0] === "/") {
    str = str.slice(1);
  }

  var newStr = str.split(" ");

  for (var i = 0; i < newStr.length; i++) {
    newStr[i] = newStr[i].slice(0, 1).toUpperCase() + newStr[i].slice(1).toLowerCase();
  }

  return newStr.join(" ");
}

function isRoute(string) {
  return /^[a-z]+$/.test(string);
}

module.exports = getDir;