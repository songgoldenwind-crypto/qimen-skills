# 排盘引擎与程序入口

默认计算来自固定提交 `8eb06d007d4a5fcc5352d9054f81469e5f023f45` 的 `atopx/qimen`。随包Go源码保留计算和测试，附JSON导出命令并清理资料出处相关注释。Python适配器只检查输入、选路、调用程序和归一化字段，不重写起局公式。软件链接与许可证见[第三方声明](../../THIRD_PARTY_NOTICES.md)。

## 配套断盘入口

需要Python 3.10+和Go 1.26。从本Skill目录运行：

```bash
bash scripts/qimen-engine/setup-atopx.sh
python3 scripts/qimen-engine/paipan.py --purpose reading --datetime 2020-04-18T14:00 --timezone Asia/Shanghai --number 7
```

`reading`是默认模式，仅允许atopx、时家、转盘、置闰和Asia/Shanghai民用时。报数可选1～9，5转坤2。程序检查八宫天／地干、门、星、神及时空是否齐备，缺项时停止。引擎或安装失败不自动换算法。

## 仅展示程序盘

其他算法没有本版配套断法，只在明确查看计算结果时调用 `--purpose chart-only`。输出 `interpretation.allowed=false`，禁止读取时间盘断法或使用报数。

```bash
python3 scripts/qimen-engine/paipan.py --purpose chart-only --datetime 2020-10-07T14:23 --style 飞盘
python3 scripts/qimen-engine/paipan.py --purpose chart-only --datetime 2020-10-07T14:23 --method 拆补
python3 scripts/qimen-engine/paipan.py --purpose chart-only --datetime 2020-10-07T14:23 --family 日家
```

Go引擎的时／日家可选转盘／飞盘、置闰／拆补；月／年家可选转盘／飞盘但不接拆补，程序固定阴遁。飞盘缺门或神的宫保留null。未与对应方法样盘逐项匹配，不由这些程序字段引入新理论。

## 可选计算依赖

刻家需要uv、Git、C++编译器，脚本安装固定版本Python依赖及固定提交源码。只能上海时区、转盘、上游 `pan_minute(2)`，不接拆补或报数。源码与环境在本机生成，不随ZIP分发。

```bash
bash scripts/qimen-engine/setup-ke.sh
python3 scripts/qimen-engine/paipan.py --purpose chart-only --family 刻家 --datetime 2026-09-21T14:23:30 --timezone Asia/Shanghai
```

提交为 `e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a`，该代码每10分钟换分柱；只检查运行、指定八宫字段和时间边界。本版不提供专用刻家断法。旬空输出来自时柱和分柱，分别标hour和minute，日空为null。

旧时家计算固定PyPI `kinqimen==0.0.6.6`，只作程序盘对照，不接报数，不改计算。它的局数及神名与默认样例有差异。

```bash
bash scripts/qimen-engine/setup.sh
scripts/qimen-engine/.venv/bin/python scripts/qimen-engine/paipan.py --purpose chart-only --engine kinqimen --datetime 2020-04-18T14:00 --method 拆补
```

## 时间与输出

输入当地墙钟时间，IANA时区另填。断盘只核上海民用时，不校真太阳时；其他时区参数只保留标记，不能视作已处理跨区交节。秒数用于记录，不承诺秒级换盘。

JSON保留原始raw、palaces、engine及版本、route、validation_scope和interpretation。`interpretation.allowed`是方法与字段的使用标记，不是准确率。默认时家空亡来自时柱；day_void与日空null不表示无日空。值符值使分别保留原宫、落宫和盘面宫。软件原值不因断语改变。

## 测试

```bash
go -C scripts/qimen-engine/atopx-src test ./...
python3 -m unittest discover -s scripts/tests -v
python3 -m unittest discover -s scripts/qimen-engine/tests -v
```

可选依赖安装后，用 `.venv/bin/python` 运行排盘测试才能检查旧引擎；未安装时对应测试跳过。固定样例与规则范围见[验证边界](../../references/10-validation-boundaries.md)。
