# 奇门排盘引擎

`paipan.py`只调用和归一化外部引擎，不自行计算局数或九宫。默认时家阳盘转盘置闰使用固定版本Go引擎，源码随Skill交付，版本与许可见[第三方声明](../../THIRD_PARTY_NOTICES.md)。副本仅清理资料出处相关注释，计算语句保留。

## 默认构建

需要Python 3.10或更新版本、Go 1.26工具链。首次构建可能下载Go工具链；生成的本机二进制不随包分发。从Skill目录运行：

```bash
bash scripts/qimen-engine/setup-atopx.sh
python3 scripts/qimen-engine/paipan.py --datetime 2020-04-18T14:00 --timezone Asia/Shanghai --number 7
python3 scripts/qimen-engine/paipan.py --family 时家 --datetime 2020-10-07T14:23 --style 飞盘
python3 scripts/qimen-engine/paipan.py --family 日家 --datetime 2020-10-07T14:23
```

支持时、刻、日、月、年家。飞盘仅用于Go引擎时／日／月／年家；时／日家可选置闰或拆补，月／年家不使用这一选项。默认仅核固定时家转盘置闰样例的指定字段；其他模式保留程序盘及未核范围，见[样盘检查](../../references/10-validation-boundaries.md)和[排盘路由](../../references/09-chart-routing.md)。

## 可选10分钟刻家

明确要求刻家才安装。需uv、Git、C++编译器及本地Python依赖；固定提交源码与虚拟环境在本机生成，不包含在交付ZIP中。

```bash
bash scripts/qimen-engine/setup-ke.sh
python3 scripts/qimen-engine/paipan.py --family 刻家 --datetime 2026-09-21T14:23:30 --timezone Asia/Shanghai
```

刻家仅限转盘、`pan_minute(2)`、上海时区，不接报数或拆补。上游虽标置闰，排局函数并不接收置闰／拆补参数；二遁按时支、三元按时柱，故为独立算法。已核10分钟边界、多日运行与八宫字段，专用断法未核。上游旬空中的日空来自时柱、时空来自分柱，归一化标为 `hour`／`minute`，真正日空保留 `null`。

## 可选旧时家对照

仅明确要求旧引擎比较时安装调用，不因默认依赖缺失自动换引擎。

```bash
bash scripts/qimen-engine/setup.sh
scripts/qimen-engine/.venv/bin/python scripts/qimen-engine/paipan.py --engine kinqimen --family 时家 --datetime 2020-04-18T14:00 --method 拆补
```

旧时家固定PyPI 0.0.6.6。适配器只修正绝对导入路径，不改计算。该版刻家对部分日期报错，不能代替单独固定的刻家提交。不同引擎的局数、星神名称和起局差异分别保留，不能拼盘。

## 输出与测试

`--datetime`是当地民用时间，可记录秒；`--timezone`为IANA时区，不校正真太阳时。秒数记录与换盘粒度不同。输出保留 `engine`、版本、`route`、原始 `raw`、九宫 `palaces`、旬空和 `validation_scope`。飞盘缺门或缺神字段为 `null`，不是吉凶判断。引擎失败或字段不完整直接报错，不补造盘面。

```bash
go -C scripts/qimen-engine/atopx-src test ./...
python3 -m unittest discover -s scripts/tests -v
python3 -m unittest discover -s scripts/qimen-engine/tests -v
```

未安装可选依赖时，对应测试明确跳过；不得据跳过声称已验证可选盘。构建检查与固定样例仍须通过。
