// 将文件夹的名称转换为标题，后续作为组件引入
// 如文件夹login --->   import Login 
function titleCase(str: string) {
    if (str[0] === "/") {
        str = str.slice(1);
    }
    var newStr = str.split(" ");
    for (var i = 0; i < newStr.length; i++) {
        newStr[i] =
            newStr[i].slice(0, 1).toUpperCase() +
            newStr[i].slice(1).toLowerCase();
    }
    return newStr.join(" ");
}

// 校验是否是路由，规则很简单，只需要a-z的组成即可
function isRoute(string: string) {
    return /^[a-z]+$/.test(string);
}

// 将所有的反斜杠都转换为斜杆，第一个斜杆替换成./
function replaceDiagonal(str: string) {
    const url = str.replace(/\\/g, () => "/");
    return url[0] === "." || url[0] === "/" ? url : "./" + url;
}

export { titleCase, isRoute, replaceDiagonal };
