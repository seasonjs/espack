//只可以使用此种注释格式
module espack_case
version 0.0.1
main index.js
//registry registry.npmjs.org //设置mod拉取的路径
target ES5 //ESNext ES5 ES2015 ES2016 ES2017 ES2018 ES2019 ES2020 ES2021  设置要编译成的语言版本

require (  //不允许进行版本提升，因此不需要lock文件，此处格式必须换行
    typescript 4.4.3
    react 17.0.2
    react-dom 17.0.2
    @types/react 17.0.24
    @types/react-dom 17.0.9
)

