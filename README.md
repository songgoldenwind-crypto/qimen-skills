# 奇门遁甲 Skill

独立开源的 `qimen-divination`，包括方法规则、排盘适配器、固定版本的默认排盘源码及测试。使用时直接读取 Skill 内的规则，不依赖六爻、梅花、私人资料目录或在线知识库。

仓库：[songgoldenwind-crypto/qimen-skills](https://github.com/songgoldenwind-crypto/qimen-skills)。可导入的 Skill ZIP 见[最新发行版](https://github.com/songgoldenwind-crypto/qimen-skills/releases/latest)。

## 能做什么

| 任务 | 输入与使用方法 |
| --- | --- |
| 一事一问、全盘主客、事业、考试、婚恋、财贸、出行、消息、失物 | 说明具体事项、对象、期限、进展、起念当地时间与时区；默认时家阳盘转盘置闰。 |
| 报数锁单宫 | 在看盘前报1～9并指定单宫法；5转坤2。评分是方法内部指标。 |
| 出生盘与家人领域 | 提供真实出生日期、时间、出生地和时区；使用出生盘的领域及房份取象。当前不提供八字起运或阴盘大运流年算法。 |
| 手机号造宫 | 提供号码末六位，逐位映射宫、神、星、门、天干、地干；无需时间盘。 |
| 住宅与选房 | 房屋中心、实测朝向和平面条件；候选房报数比较与已居住宅方位分析分别使用。 |
| 择时择方 | 提供实际事项、地点、方向和候选办事时刻，每个候选时刻分别排盘。 |
| 已有盘复核 | 提供原盘及家法、盘式、起局法和时间，保留原值核对。 |
| 可选刻家、飞盘及日／月／年家 | 明确指定入口，保留各自算法及验证范围。10分钟刻家另需安装依赖，当前缺少相容的专用断法。 |

阴盘／阳盘、阴遁／阳遁、转盘／飞盘、置闰／拆补分别说明，不能拼成一组选项。阴盘暂无经样盘核验的自动起局入口。时间可记到秒，换盘粒度由所选算法决定。固定样例校验不等于现实预测准确率。

## 导入与使用

实际 Skill 位于 `skills/qimen-divination/`；导入时选择整个目录，或解压 `dist/qimen-divination.zip` 得到同名目录，保留其 `SKILL.md`、`references/`、`scripts/`、`agents/` 和许可证。环境须支持本地Python执行；默认引擎首次构建还需Go 1.26工具链。

导入后可输入：

```text
$qimen-divination 用当前时间按时家奇门分析这次求职，期限是本周收到面试通知。
$qimen-divination 我预先报数7，用报数单宫法分析这件事。
$qimen-divination 用末六位377287做手机号造宫分析。
```

开发目录中的首次构建和排盘：

```bash
bash skills/qimen-divination/scripts/qimen-engine/setup-atopx.sh
python3 skills/qimen-divination/scripts/qimen-engine/paipan.py --datetime 2026-09-21T14:20 --timezone Asia/Shanghai
```

引擎安装、可选依赖及算法范围见[引擎说明](skills/qimen-divination/scripts/qimen-engine/README.md)，能力入口见[SKILL.md](skills/qimen-divination/SKILL.md)。

## 检查与打包

```bash
python3 scripts/check_release.py
python3 scripts/package_skill.py
```

检查覆盖目录结构、引用链接、文档和代码中的出处标识，以及交付ZIP的文件内容。运行环境、编译产物、缓存、Git历史和原始文献不进入ZIP。方法内容不携带资料书名、作者、页码、图号、编号或私有路径，推断输出也只说明方法与条件。

本项目采用MIT开源许可，允许使用、修改及分发，见[LICENSE](LICENSE)；第三方排盘软件独立适用其许可证，见[第三方声明](THIRD_PARTY_NOTICES.md)。
