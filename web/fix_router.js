const fs = require('fs');
const path = 'd:\\goprojects\\triad\\web\\src\\router\\index.js';
let content = fs.readFileSync(path, 'utf8');
content = content.replace(/=&gt;/g, '=&gt;');
fs.writeFileSync(path, content, 'utf8');
console.log('修复完成！');
