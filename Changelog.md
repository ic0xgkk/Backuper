# Changelog

## v1.2

* 修正了文件描述符的关闭逻辑，解决了之前文件损坏的问题
* 新增了日志等级设置
* 合并了日志模块，现在共用日志句柄但根据任务进度不同日志的Writer不同
* 不再支持Windows
* 增加了自动删除选项
* 修改文件压缩记录为Hook模式
