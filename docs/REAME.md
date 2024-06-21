# App兼容性工具deepin-app-analyze
    App应用中的库ABI是`Application Binary Interface`的缩写，以二进制形式发布动态库时，需要检查ABI兼容性保证在系统运行时的稳定，由此而产生检查工具deepin-app-analyze

## 内容
deepin-app-analyze主要用于对比App的deb包中所包含的二进制库文件与系统基线对比而分析

### 网络

  - \[ \] 网络连接

> 目前是离线扫描工具不依赖网络环境.


### 更新

  - \[X\] 自动分析检查(`AutoAnalysisCheck`)
  - \[X\] 自动生成报告(`AutoGeneratReport`)

### 使用
---
    appcheck  (app检查命令常用操作)
    `deepin-app-analyze appcheck -f package`
---
    appcheck -e/-example  (app检查命令使用示例)
    `deepin-app-analyze appcheck -e test`

