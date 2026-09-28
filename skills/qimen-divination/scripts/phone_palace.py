#!/usr/bin/env python3
"""Six-digit phone-number single-palace conversion (no time chart)."""

from __future__ import annotations

import argparse
import json

PALACE = ("宫空亡", "坎1", "坤2", "震3", "巽4", "坤2（中5寄）", "乾6", "兑7", "艮8", "离9")
SPIRIT = ("神空亡", "值符", "螣蛇", "太阴", "六合", "白虎", "玄武", "九地", "九天", "值符")
STAR = ("星空亡", "天蓬", "天芮", "天冲", "天辅", "天禽", "天心", "天柱", "天任", "天英")
DOOR = ("门空亡", "休门", "死门", "伤门", "杜门", "死门（中5寄）", "开门", "惊门", "生门", "景门")
STEM = ("癸", "甲（戊）", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬")


def build_palace(number: str) -> dict[str, str]:
    if not number.isascii() or not number.isdigit() or len(number) < 6:
        raise ValueError("号码必须是至少六位的 ASCII 数字；保留前导 0")
    digits = number[-6:]
    values = [int(x) for x in digits]
    return {
        "method": "手机号末六位造宫",
        "last_six_digits": digits,
        "palace": PALACE[values[0]],
        "spirit": SPIRIT[values[1]],
        "star": STAR[values[2]],
        "door": DOOR[values[3]],
        "heaven_stem": STEM[values[4]],
        "earth_stem": STEM[values[5]],
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="将号码末六位映射为奇门单宫")
    parser.add_argument("number", help="完整号码或末六位；只接受 ASCII 数字")
    args = parser.parse_args()
    try:
        result = build_palace(args.number)
    except ValueError as exc:
        parser.error(str(exc))
    print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
