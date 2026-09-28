---
name: qimen-divination
description: 用于时家阳盘转盘置闰的报数单宫、宫位关系、具体事项、出生领域与择时分析，以及手机号末六位造宫。其他排盘算法仅展示程序盘，不提供配套断法。
---

# 奇门遁甲

本版断盘只配套 `atopx/qimen` 固定版本的时家阳盘转盘置闰，采用Asia/Shanghai民用时间。只使用下列已适配规则，不临时借入其他盘法的取用、寄宫或应期公式。说明具体方法与条件，不输出资料书名、作者、页码、图号、来源编号或本地资料路径；软件版本与必要许可证独立保留。

## 起盘与规则检查

记录事项、对象、期限、现实进展及起念时间。出生领域用真实出生时刻；择时用每个候选办事时刻。未知分钟不编造。所有脚本路径相对本Skill目录，首次构建见[引擎说明](scripts/qimen-engine/README.md)。

```bash
bash scripts/qimen-engine/setup-atopx.sh
python3 scripts/qimen-engine/paipan.py --purpose reading --datetime 2026-09-21T14:20 --timezone Asia/Shanghai
```

起盘后先核 `route.purpose`、`interpretation.allowed`、`interpretation.rule_set` 和 `validation_scope`。只有 `reading`、`allowed=true` 且规则为“时家阳盘转盘置闰”时读取断盘规则。程序拒绝非配套盘法或八宫必要字段缺失；缺依赖时说明安装要求，不换引擎补盘。

若明确只要其他算法的程序盘，按[排盘路由](references/09-chart-routing.md)调用 `--purpose chart-only`。此模式只说明原值、算法版本与计算限制，`interpretation.allowed=false`；不进入下列断盘规则。阴盘尚无自动排盘入口；阴盘或方法不明的截图只整理字段，不套用本版断法。

## 已适配规则

| 任务 | 具体入口 |
| --- | --- |
| 用户事前报1～9，选择锁单宫 | 加 `--number N`，读[报数单宫](references/01-single-palace.md)；报5转坤2。 |
| 未报数、查看多个宫或比较关系 | 读[宫位关系](references/05-full-chart.md)；使用明确选定的宫及实际字段，不代报数、不套单宫评分。 |
| 求职、考试、求财、会面、出行、消息或失物 | 读[事项解释](references/06-question-routes.md)；事项只限定本题含义，不引入其他方法的专用用神。 |
| 同宫天／地盘干组合 | 读[81组干组合](references/02-stem-patterns.md)，干序明确才查表。 |
| 宫、门、星、神 | 按需读[符号速查](references/07-symbols.md)，以实际落宫字段为准。 |
| 出生领域、家属房份、住宅或择时 | 读[特殊应用](references/03-special-applications.md)，分别核输入和条件。 |
| 手机号末六位造宫 | 单独调用 `scripts/phone_palace.py`，读[特殊应用](references/03-special-applications.md)；这是确定性数字映射，无时间盘。 |

[规则与字段](references/04-method-compatibility.md)列明每条方法所需的程序字段；[时间与盘式](references/08-chart-methods.md)说明二遁、换盘与民用时间；[样盘检查](references/10-validation-boundaries.md)限定已核范围。

## 解读与输出

先回答本题支持到哪一步，再给实际时间、规则、引擎版本和关键宫值。判断按“实际字段 → 本版规则与条件 → 事项解释 → 待核点”展开；原始值与解释分开。中宫无门星神时不补造，日空为 `null` 时不当作无空亡或另算日空。

同一时辰可以共基础盘，各问题分别固定对象、选宫和期限；同盘多次解释不算独立复验。择时只比较实际排出的候选时刻；本版不提供专门流年、大运或精确年月日应期公式。单宫分数不换算现实概率，符号不作为诊断、开奖号码、隐私事实或收益承诺。
