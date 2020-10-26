function titleCase(str) {
  if (str[0] === "/") {
    str = str.slice(1);
  }
  var newStr = str.split(" ");
  for (var i = 0; i < newStr.length; i++) {
    newStr[i] =
      newStr[i].slice(0, 1).toUpperCase() + newStr[i].slice(1).toLowerCase();
  }
  return newStr.join(" ");
}

function isRoute(string) {
  return /^[a-z]+$/.test(string);
}

function replaceDiagonal(str) {
  return str.replace(/\\/g, () => "/");
}

module.exports = { titleCase, isRoute, replaceDiagonal };
