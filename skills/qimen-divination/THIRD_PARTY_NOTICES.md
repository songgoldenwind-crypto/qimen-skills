# 第三方软件许可

自有规则及适配器采用MIT开源许可。必要的第三方版权与许可证独立保留，不涉及占断资料的出处。

- 默认Go引擎来自 `https://github.com/atopx/qimen`，固定提交 `8eb06d007d4a5fcc5352d9054f81469e5f023f45`，采用MIT许可；原许可证保留在 `scripts/qimen-engine/atopx-src/LICENSE`。副本清理资料出处相关注释，附JSON导出命令；排盘计算语句和上游测试保留，本机二进制不分发。
- 可选旧时家引擎使用PyPI `kinqimen==0.0.6.6`，安装时取得的依赖各依其自身许可证。本项目不分发这些依赖或二进制。
- 可选刻家安装时从 `https://github.com/kentang2017/kinqimen` 取得固定提交 `e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a`。该提交的README声明MIT，但缺少所链接的LICENSE文件；源码不包含在交付包中。
- 可选安装脚本另固定 `sxtwl_cpp` 提交 `7598b0601a76cfdaa9266257b1b5690720c1e2ce`、`bidict==0.23.1`及`ephem==4.1.6`。它们在本机安装，分别适用取得时的许可证。

这些软件许可不限制自有规则的编排，也不改变用户依法取得的第三方软件权利。
