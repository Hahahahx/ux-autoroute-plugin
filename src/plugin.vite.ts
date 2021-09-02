import * as path from "path";
import * as process from "child_process";

export function AutoRouterVitePlugin(params: UxPlugin.PluginParams) {
  return {
    name: "RouterPlugin", // 名称用于警告和错误展示
    buildStart() {
      buildRouter(params);
    },
    handleHotUpdate() {
      buildRouter(params);
    },
  };
}

function buildRouter(params: UxPlugin.PluginParams) {
  const parse_exe = path.join(__dirname, "..", "parse.exe");
  process.exec(
    `${parse_exe} parse -o ${params.output} -r ${params.pagePath} -n ${
      params.filename
    } ${!!params.defaultLazyImport ? "-i" : ""}`,
    (err: any, stdout: any, stderr: any) => {
      // console.log(stdout);
    }
  );
}

// module.exports = { AutoRouterVitePlugin };
