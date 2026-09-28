# 奇门遁甲 Skill

独立开源的 `qimen-divination`。**本版断盘规则只配套时家阳盘转盘置闰，排盘使用固定版本的 `atopx/qimen` Go代码。** 规则只保留当前盘面能支持的报数、宫位解释、出生领域和择时方法。手机号造宫是另一个确定性数字映射入口。

其他引擎、拆补、飞盘、刻家或年／月／日家可以显式导出程序盘，但本版不为它们提供配套断法。**目前没有现代阴盘奇门自动起局、终身大运流年或专门精确应期算法。**

仓库：[songgoldenwind-crypto/qimen-skills](https://github.com/songgoldenwind-crypto/qimen-skills)；可导入ZIP见[最新发行版](https://github.com/songgoldenwind-crypto/qimen-skills/releases/latest)。

## 排盘代码从哪里来

| 组件 | 实际工作与代码位置 | 固定版本 |
| --- | --- | --- |
| 默认排盘计算 | [atopx/qimen](https://github.com/atopx/qimen/tree/8eb06d007d4a5fcc5352d9054f81469e5f023f45)；Go源码随Skill放在 `scripts/qimen-engine/atopx-src/`，计算历法、干支、局数、二遁、门星神与九宫。 | 提交 `8eb06d007d4a5fcc5352d9054f81469e5f023f45`；MIT。 |
| Go JSON导出命令 | `atopx-src/cmd/qimen-json/main.go` 调用上述引擎，将原盘字段导出JSON；本地构建 `.bin/atopx-qimen`。 | 随本项目版本维护。 |
| Python适配器 | `scripts/qimen-engine/paipan.py` 负责输入校验、固定路由、进程调用、字段归一化与规则支持检查。局数、门星神不由Python或LLM重算。 | 随本项目版本维护。 |
| 可选旧时家计算 | PyPI `kinqimen==0.0.6.6`，只作程序盘对照；安装时适配绝对导入，不改原计算。 | `0.0.6.6`；局数与神名在固定样例中有差异。 |
| 可选10分钟刻家计算 | [kentang2017/kinqimen](https://github.com/kentang2017/kinqimen/tree/e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a)，调用该提交的 `pan_minute(2)`；源码安装到本机 `.ke-src/`，不随ZIP分发。 | 提交 `e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a`。 |
| 手机号造宫 | `scripts/phone_palace.py`，取末六位逐位映射宫、神、星、门与两干。 | 项目内确定性映射，不是时间起局。 |

随包Go源码保留上游排盘计算和测试，附JSON导出边界，并清理资料出处相关注释。上游版本号、软件仓库链接及必要版权许可证保留；占断规则不携带资料书名、资料作者、页码或出处编号。软件许可及可选依赖版本见[第三方声明](THIRD_PARTY_NOTICES.md)。

## 哪些排盘可以解读

| 引擎／算法 | 可以排程序盘 | 可以调用本版断法 |
| --- | --- | --- |
| 默认atopx时家／阳盘／转盘／置闰，Asia/Shanghai民用时 | 可以；默认入口。 | 可以，仍限于下述样例验证与本版规则。 |
| atopx时家拆补、飞盘、日／月／年家 | 可显式选 `--purpose chart-only`。 | 不可以。 |
| 旧kinqimen时家 | 需安装可选依赖，只允许chart-only。 | 不可以。 |
| 固定提交10分钟刻家 | 需安装可选依赖，限上海时区，只允许chart-only。 | 不可以；当前没有同算法配套断法。 |
| 现代阴盘 | 尚未接入自动排盘。 | 不提供。 |

`--purpose reading` 是默认模式。非配套算法或时区在计算前被拒绝；八宫必要字段缺失时也停止。只有 `interpretation.allowed=true` 才能读取本版时间盘规则。`chart-only` 输出强制标为 `allowed=false`，仅展示字段，不拿默认断法继续解释。

不同阴阳遁可以出现在当前阳盘算法中。换成阴遁不是切换到现代阴盘；其他算法输出一张盘也不表示已完成理论配套。

## 当前保留的能力

| 任务 | 输入与做法 |
| --- | --- |
| 报数单宫 | 起念时间、具体问题、事前报1～9并选择此法；5转坤2。评分是方法内部指标。 |
| 宫位与具体事项 | 问题、对象、期限与实际进展；读取事类相关宫或事前明确的宫，比较实际符号及五行关系。 |
| 同宫干组合 | 明确的天盘干和地盘干，按顺序查81组名称及本题含义；不补入未定义的附格。 |
| 出生领域与家属房份 | 真实出生时间、出生地、时区和已确认家庭关系；只读固定领域，不推大运流年。 |
| 住宅与选房 | 候选房预报数比较，或已居住宅的房屋中心、实测方向与现场情况。 |
| 择时择方 | 每个候选办事时刻分别排默认盘，按休／生／开门、宫位和排除条件筛选。 |
| 手机号末六位造宫 | 六个ASCII数字，独立逐位映射；不与时间盘空亡或报数法混用。 |

未配套的专用取用、寄宫、阴盘布局和应期理论已从规则正文移除，不再以“另法对照”参与断语。固定出生盘的领域解释不等于完整终身命理；择时比较实际候选时间也不等于精确推算事件日期。

## 时间与输出口径

- 断盘采用Asia/Shanghai当地民用时间；不校正真太阳时。其他时区的交节及夏令时未适配，仅能显式展示未核程序盘。
- 输入可保留秒数，常规时家以时辰为基础；秒数不是秒级换盘能力，交节与换日等边界按固定引擎处理。
- 输出保留 `engine/engine_revision`、`route`、`interpretation`、`validation_scope`、`raw` 和 `palaces`。
- 默认时家旬空来自时柱，`hour_void` 是时空；`day_void` 及 `void_branches.day` 为 `null`，本版不补算日空。
- 中5的门星神可能为 `null`；保留原值，不当作第九个完整评分宫。值符值使的原宫、落宫与盘面宫分别保留；“禽芮”不覆盖成“天芮”。

## 验证做到哪一步

默认引擎已核三组固定时家转盘置闰样例的局数及指定宫门、神、天／地干，另核两组兑7门变化。部分星名仍存在差异，未宣称所有日期和全宫逐项相同。具体时间与字段见[验证边界](skills/qimen-divination/references/10-validation-boundaries.md)。

上游Go测试验证算法内部行为；本项目测试验证默认路由、方法拒绝、字段完整性、程序盘标记、固定样例与数字映射。可选依赖未安装时相应测试会跳过，跳过不算通过。**结构正确、字段匹配和样例测试均不等于现实预测准确率。**

## 安装与使用

导入 `skills/qimen-divination/` 整个目录，或从发行版下载 `qimen-divination.zip`，解压后导入同名目录。需要本地Python 3.10+；默认引擎首次构建需要Go 1.26工具链，可能下载工具链。编译产物与运行环境不包含在ZIP中。

在仓库根目录首次构建和排盘：

```bash
bash skills/qimen-divination/scripts/qimen-engine/setup-atopx.sh
python3 skills/qimen-divination/scripts/qimen-engine/paipan.py --purpose reading --datetime 2026-09-21T14:20 --timezone Asia/Shanghai
```

导入后可输入：

```text
$qimen-divination 用当前时间分析这次求职，期限是本周收到面试通知。
$qimen-divination 我预先报数7，用报数单宫法分析这件事。
$qimen-divination 用末六位377287做手机号造宫分析。
```

明确只看其他算法程序盘，例如飞盘：

```bash
python3 skills/qimen-divination/scripts/qimen-engine/paipan.py --purpose chart-only --style 飞盘 --datetime 2026-09-21T14:20 --timezone Asia/Shanghai
```

完整参数、可选依赖与输出限制见[引擎说明](skills/qimen-divination/scripts/qimen-engine/README.md)，能力入口见[SKILL.md](skills/qimen-divination/SKILL.md)。

## 检查与打包

```bash
python3 scripts/check_release.py
python3 scripts/package_skill.py
```

发布检查覆盖结构、内部链接、资料出处标识及ZIP内容一致性；原始文献、个人路径、运行环境、缓存和Git历史不进入交付ZIP。

自有规则和适配器采用[MIT许可证](LICENSE)，第三方软件按各自许可证使用。
