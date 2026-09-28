#!/usr/bin/env python3
"""Label and normalize pinned Qimen engines; never calculate a chart here."""

from __future__ import annotations

import argparse
import importlib
import importlib.metadata
import json
import subprocess
import sys
from datetime import datetime, timedelta
from pathlib import Path
from zoneinfo import ZoneInfo, ZoneInfoNotFoundError


ENGINE_VERSION = "0.0.6.6"
ATOPX_REVISION = "8eb06d007d4a5fcc5352d9054f81469e5f023f45"
KE_REVISION = "e6680ac4ca0b0da5ce3fe637e05f9fc32066ec5a"
ENGINE_METHODS = {"置闰": (2, "置閏"), "拆补": (1, "拆補")}
FAMILIES = {"时家": "time", "日家": "day", "月家": "month", "年家": "year"}
ALL_FAMILIES = ("时家", "刻家", "日家", "月家", "年家")
STYLES = {"转盘": "rotate", "飞盘": "fly"}
JU_DIGITS = "零一二三四五六七八九"
PALACES = {1: "坎", 2: "坤", 3: "震", 4: "巽", 5: "中", 6: "乾", 7: "兌", 8: "艮", 9: "離"}
BRANCH_PALACE = dict(zip("子丑寅卯辰巳午未申酉戌亥", (1, 8, 8, 3, 4, 4, 9, 2, 2, 7, 6, 6)))
SPIRIT_NAMES = {"符": "值符", "蛇": "螣蛇", "陰": "太陰", "合": "六合", "虎": "白虎", "玄": "玄武", "勾": "勾陈", "雀": "朱雀", "地": "九地", "天": "九天"}


def load_engine():
    """Work around the upstream wheel's absolute `import config` without editing it."""
    try:
        installed = importlib.metadata.version("kinqimen")
    except importlib.metadata.PackageNotFoundError as exc:
        raise RuntimeError("缺少 kinqimen==0.0.6.6；见 scripts/qimen-engine/README.md") from exc
    if installed != ENGINE_VERSION:
        raise RuntimeError(f"需要 kinqimen=={ENGINE_VERSION}，当前为 {installed}；其他版本未通过固定样例校验")
    try:
        package = importlib.import_module("kinqimen")
        package_path = str(Path(package.__file__).resolve().parent)
        if package_path not in sys.path:
            sys.path.insert(0, package_path)
        return importlib.import_module("kinqimen.kinqimen").Qimen
    except (ImportError, OSError) as exc:
        raise RuntimeError(f"kinqimen 依赖缺失或载入失败：{exc}；见 scripts/qimen-engine/README.md") from exc


def parse_local_time(value: str, timezone: str) -> datetime:
    try:
        local = datetime.fromisoformat(value)
    except ValueError as exc:
        raise ValueError("--datetime 应为 YYYY-MM-DDTHH:MM[:SS]") from exc
    if local.tzinfo is not None:
        raise ValueError("--datetime 请填当地墙钟时间，时区另用 --timezone 指定")
    try:
        zone = ZoneInfo(timezone)
    except ZoneInfoNotFoundError as exc:
        raise ValueError(f"未知 IANA 时区：{timezone}") from exc
    return local.replace(tzinfo=zone)


def resolve_route(family: str, style: str, method: str, number: int | None, engine: str) -> str:
    """Select only a compatible engine; 阴/阳遁 is derived by that engine."""
    if family not in ALL_FAMILIES or style not in STYLES or method not in ENGINE_METHODS:
        raise ValueError("不支持的家法、盘式或起局方法")
    if engine not in {"auto", "atopx", "kinqimen", "kinqimen-ke"}:
        raise ValueError("不支持的排盘引擎")
    if family == "刻家":
        if engine not in {"auto", "kinqimen-ke"}:
            raise ValueError("刻家仅使用固定提交的 KinQiMen 刻家算法")
        if style != "转盘" or method != "置闰" or number is not None:
            raise ValueError("刻家只开放上游 pan_minute(2) 转盘，不套用时家报数或拆补")
        return "kinqimen-ke"
    if engine == "kinqimen-ke":
        raise ValueError("KinQiMen 刻家不能冒充时家、日家、月家或年家")
    if engine == "kinqimen" and (family != "时家" or style != "转盘"):
        raise ValueError("旧 KinQiMen 只作时家转盘来源对照")
    if family in {"月家", "年家"} and method != "置闰":
        raise ValueError("月家和年家无置闰/拆补选择")
    if number is not None and (family != "时家" or style != "转盘"):
        raise ValueError("报数锁宫仅用于时家转盘")
    return "atopx" if engine == "auto" else engine


def build_chart(local: datetime, number: int | None, method: str = "置闰") -> dict:
    if method not in ENGINE_METHODS:
        raise ValueError(f"未知排盘方式：{method}")
    method_number, expected_method = ENGINE_METHODS[method]
    Qimen = load_engine()
    raw = Qimen(local.year, local.month, local.day, local.hour, local.minute).pan(method_number)
    if raw.get("排盤方式") != expected_method:
        raise RuntimeError(f"上游未返回{method}盘，已停止")
    required = ("天盤", "地盤", "門", "星", "神", "旬空", "排局")
    if any(key not in raw for key in required):
        raise RuntimeError("上游盘面缺少必要字段，已停止")
    hour_void_branches = raw["旬空"].get("時空", "")
    day_void_branches = raw["旬空"].get("日空", "")
    hour_void_palaces = {BRANCH_PALACE[b] for b in hour_void_branches if b in BRANCH_PALACE}
    day_void_palaces = {BRANCH_PALACE[b] for b in day_void_branches if b in BRANCH_PALACE}
    palaces = {}
    for palace_number, name in PALACES.items():
        sky = raw["天盤"].get(name)
        earth = raw["地盤"].get(name)
        star = raw["星"].get(name)
        door = raw["門"].get(name)
        spirit = raw["神"].get(name)
        palaces[str(palace_number)] = {
            "name": name,
            "sky_stem": sky,
            "earth_stem": earth,
            "stem_pair": f"{sky}+{earth}" if sky and earth else None,
            "door": door,
            "star": star,
            "spirit": SPIRIT_NAMES.get(spirit, spirit),
            "hour_void": palace_number in hour_void_palaces,
            "day_void": palace_number in day_void_palaces,
        }
    locked = 2 if number == 5 else number
    return {
        "engine": "kinqimen",
        "engine_version": ENGINE_VERSION,
        "validation_scope": "旧时家对照；固定样例存在起局和神名差异",
        "method": f"时家奇门/阳盘/转盘/{method}",
        "engine_method_number": method_number,
        "source_time": local.isoformat(timespec="minutes"),
        "timezone": local.tzinfo.key,
        "reported_number": number,
        "locked_palace": locked,
        "five_redirected_to_two": number == 5,
        "void_branches": {"day": day_void_branches, "hour": hour_void_branches},
        "raw": raw,
        "palaces": palaces,
    }


def build_ke_chart(local: datetime, method: str = "置闰") -> dict:
    """Call KinQiMen's pinned ten-minute chart in an isolated Python process."""
    if method != "置闰":
        raise ValueError("刻家当前只开放已检查的 KinQiMen 置闰入口")
    if getattr(local.tzinfo, "key", None) != "Asia/Shanghai":
        raise ValueError("KinQiMen 刻家目前仅按 Asia/Shanghai 当地民用时核对；其他时区不可直接套用")
    engine_dir = Path(__file__).resolve().parent
    runtime = engine_dir / ".venv" / "bin" / "python"
    if not runtime.is_file():
        raise RuntimeError("刻家依赖尚未安装；先运行 scripts/qimen-engine/setup-ke.sh")
    result = subprocess.run(
        [str(runtime), str(engine_dir / "ke_paipan.py"),
         local.strftime("%Y-%m-%dT%H:%M:%S"), "2"],
        capture_output=True, text=True, check=False,
    )
    if result.returncode:
        raise RuntimeError(f"刻家引擎排盘失败：{result.stderr.strip() or result.returncode}")
    try:
        raw = json.loads(result.stdout)
        if raw["排盤方式"] != "置閏" or not all(
            isinstance(raw[key], dict) for key in ("天盤", "地盤", "門", "星", "神", "旬空")
        ):
            raise ValueError("盘式或九宫字段不完整")
        if not raw["排局"] or not raw["干支"]:
            raise ValueError("缺少排局或干支")
        for key in ("天盤", "地盤", "門", "星", "神"):
            if not all(name in raw[key] for name in PALACES.values() if name != "中"):
                raise ValueError(f"{key}缺少八宫")
    except (json.JSONDecodeError, KeyError, TypeError, ValueError) as exc:
        raise RuntimeError(f"刻家引擎返回的盘面不完整：{exc}") from exc
    # Upstream hourkong_minutekong() puts the hour pillar under 日空 and
    # the ten-minute pillar under 時空. Keep raw keys but label their origin.
    hour_void = raw["旬空"].get("日空", "")
    minute_void = raw["旬空"].get("時空", "")
    hour_void_palaces = {BRANCH_PALACE[b] for b in hour_void if b in BRANCH_PALACE}
    minute_void_palaces = {BRANCH_PALACE[b] for b in minute_void if b in BRANCH_PALACE}
    palaces = {}
    for palace_number, name in PALACES.items():
        sky = raw["天盤"].get(name)
        earth = raw["地盤"].get(name)
        spirit = raw["神"].get(name)
        palaces[str(palace_number)] = {
            "name": name,
            "sky_stem": sky,
            "earth_stem": earth,
            "stem_pair": f"{sky}+{earth}" if sky and earth else None,
            "door": raw["門"].get(name),
            "star": raw["星"].get(name),
            "spirit": SPIRIT_NAMES.get(spirit, spirit),
            "hour_void": palace_number in hour_void_palaces if hour_void else None,
            "minute_void": palace_number in minute_void_palaces if minute_void else None,
        }
    start = local.replace(minute=(local.minute // 10) * 10, second=0, microsecond=0)
    end = start + timedelta(minutes=10)
    return {
        "engine": "kentang2017/kinqimen",
        "engine_revision": KE_REVISION,
        "method": "刻家奇门/10分钟/转盘/上游标注置闰",
        "source_time": local.isoformat(timespec="seconds"),
        "timezone": local.tzinfo.key,
        "time_basis": "当地民用时；未校正真太阳时",
        "chart_granularity_minutes": 10,
        "ke_interval": {"start_inclusive": start.isoformat(timespec="seconds"),
                        "end_exclusive": end.isoformat(timespec="seconds")},
        "validation_scope": "刻家程序盘；已检查10分钟边界及多日运行，同算法专用断法尚未核验",
        "method_caveat": "上游刻家阴阳遁按时支分段、三元按时柱取得；排局函数不接收置闰/拆补选项。不能套用时家节气二遁、时家校验结论或报数法",
        "reported_number": None,
        "locked_palace": None,
        "void_branches": {"hour": hour_void, "minute": minute_void,
                          "day": None, "upstream_raw": raw["旬空"]},
        "raw": raw,
        "palaces": palaces,
    }


def build_atopx_chart(
    local: datetime,
    number: int | None,
    method: str = "置闰",
    family: str = "时家",
    style: str = "转盘",
) -> dict:
    """Read a chart from the pinned MIT-licensed Go engine, without merging sources."""
    if family not in FAMILIES or style not in STYLES or method not in ENGINE_METHODS:
        raise ValueError("不支持的家法、盘式或起局方法")
    if number is not None and (family != "时家" or style != "转盘"):
        raise ValueError("报数锁宫仅用于时家转盘")
    if family in {"月家", "年家"} and method != "置闰":
        raise ValueError("月家和年家不使用置闰/拆补选项")
    binary = Path(__file__).resolve().parent / ".bin" / "atopx-qimen"
    if not binary.is_file():
        raise RuntimeError("尚未构建 atopx 引擎；先运行 scripts/qimen-engine/setup-atopx.sh")
    result = subprocess.run(
        [
            str(binary), "--datetime", local.strftime("%Y-%m-%dT%H:%M:%S"),
            "--family", FAMILIES[family], "--style", STYLES[style],
            "--method", "zhirun" if method == "置闰" else "chaibu",
        ],
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode:
        raise RuntimeError(f"atopx 引擎排盘失败：{result.stderr.strip() or result.returncode}")
    try:
        chart = json.loads(result.stdout)
        raw, palaces = chart["raw"], chart["palaces"]
        ju = raw["局数"]
        if not (1 <= ju <= 9) or set(palaces) != {str(i) for i in range(1, 10)}:
            raise ValueError("缺少九宫或局数不合法")
        raw["排局"] = f"{raw['阴阳遁']}遁{JU_DIGITS[ju]}局{raw['三元']}元"
    except (ValueError, KeyError, TypeError) as exc:
        raise RuntimeError(f"atopx 引擎返回的盘面不完整：{exc}") from exc
    locked = 2 if number == 5 else number
    rule_label = method if family in {"时家", "日家"} else "不适用"
    validated_case = family == "时家" and style == "转盘" and method == "置闰"
    method_parts = [f"{family}奇门"]
    if family == "时家" and style == "转盘":
        method_parts.append("阳盘")
    method_parts.extend((style, rule_label))
    return {
        "engine": "atopx/qimen",
        "engine_revision": ATOPX_REVISION,
        "method": "/".join(method_parts),
        "source_time": local.isoformat(timespec="seconds"),
        "timezone": local.tzinfo.key,
        "time_basis": "当地民用时；未校正真太阳时",
        "validation_scope": (
            "固定时家转盘置闰样例的指定字段及兑7门变化已核；部分星名有差异，非全盘逐项等同"
            if validated_case else "引擎可排此盘式；尚未与对应独立方法样盘逐项核对"
        ),
        "reported_number": number,
        "locked_palace": locked,
        "five_redirected_to_two": number == 5,
        "void_branches": {"lead": "".join(raw["旬空"]), "day": None},
        "raw": raw,
        "palaces": palaces,
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="来源分离的奇门排盘：默认指定样例字段已核的时家转盘置闰")
    parser.add_argument("--datetime", required=True, help="当地墙钟时间，如 2020-10-07T14:23")
    parser.add_argument("--timezone", default="Asia/Shanghai", help="IANA 时区；默认 Asia/Shanghai")
    parser.add_argument("--engine", choices=("auto", "atopx", "kinqimen", "kinqimen-ke"), default="auto", help="auto 依所选家法固定来源；不能按吉凶换引擎")
    parser.add_argument("--family", choices=ALL_FAMILIES, default="时家", help="默认时家；10分钟刻家须显式选择")
    parser.add_argument("--style", choices=STYLES, default="转盘", help="飞盘只在 atopx 中可用")
    parser.add_argument("--number", type=int, choices=range(1, 10), metavar="1..9", help="问测人预先报出的数")
    parser.add_argument("--method", choices=ENGINE_METHODS, default="置闰", help="时家/日家选置闰或拆补；不能按结果自动切换")
    args = parser.parse_args()
    try:
        local = parse_local_time(args.datetime, args.timezone)
        engine = resolve_route(args.family, args.style, args.method, args.number, args.engine)
        if engine == "kinqimen-ke":
            chart = build_ke_chart(local, args.method)
        elif engine == "kinqimen":
            chart = build_chart(local, args.number, args.method)
        else:
            chart = build_atopx_chart(local, args.number, args.method, args.family, args.style)
    except (ValueError, RuntimeError) as exc:
        parser.exit(2, f"error: {exc}\n")
    raw_dun = chart["raw"].get("阴阳遁", chart["raw"].get("排局", ""))
    dun = "阴遁" if raw_dun.startswith(("阴", "陰")) else "阳遁" if raw_dun.startswith(("阳", "陽")) else None
    if dun is None:
        parser.exit(2, "error: 引擎未提供可识别的阴阳遁，已停止\n")
    dun_basis = (
        "KinQiMen 刻家源码按时支子至巳/午至亥分阳遁/阴遁；不是时家节气二遁"
        if args.family == "刻家" else
        "atopx 年/月家程序固定阴遁；尚未与独立样盘核验"
        if args.family in {"月家", "年家"} else
        "选定引擎按时家/日家用局节气计算；不是用户另选的阴盘/阳盘"
    )
    chart["route"] = {
        "family": args.family,
        "style": args.style,
        "rule_requested": args.method if args.family in {"时家", "日家"} else ("上游 pan_minute(2)" if args.family == "刻家" else "不适用"),
        "engine": engine,
        "yin_yang_dun": dun,
        "yin_yang_source": dun_basis,
    }
    print(json.dumps(chart, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
