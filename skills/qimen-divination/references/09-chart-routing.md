# 排盘路由：先核规则是否可用

| 任务 | 实际调用 | 后续动作 |
| --- | --- | --- |
| 起念问事或普通宫位分析 | `paipan.py --purpose reading --datetime ... --timezone Asia/Shanghai` | 默认atopx时家阳盘转盘置闰；检查 `interpretation.allowed=true` 后读本版规则。 |
| 已预报1～9，选择锁宫 | 同上加 `--number N` | 读取锁宫；5转坤2。 |
| 出生领域 | 同一reading入口，时间填真实出生时刻 | 只读领域与已确认房份，不延伸成流年大运算法。 |
| 择时择方 | 每个候选办事时刻分别调用同一reading入口 | 按目标门宫和排除条件比较。 |
| 手机号末六位造宫 | `phone_palace.py 377287` | 独立数字映射，不调用时间盘或报数参数替代。 |
| 明确只要其他引擎／算法的程序盘 | `paipan.py --purpose chart-only ...` | 仅展示原值、程序版本和未核范围，`interpretation.allowed=false`；不调用断盘规则。 |
| 用户已有截图 | 保留输入原值并核方法、时间和字段 | 仅当默认规则、版本和必要字段可确认一致时读本版规则；否则只整理字段。 |

默认 `--purpose` 是 `reading`。非默认引擎、家法、飞转或起局法会在排盘前被拒绝，避免显示一张无配套断法的盘却继续断事。其他算法只在明确要程序盘时使用 `chart-only`，技术参数见[引擎说明](../scripts/qimen-engine/README.md)。

输出核 `route.family/style/rule_requested/engine/purpose/yin_yang_dun`、`interpretation` 和 `validation_scope`。阴阳遁由程序决定；缺字段、缺依赖或无法确认方法时不凭模型补盘。
