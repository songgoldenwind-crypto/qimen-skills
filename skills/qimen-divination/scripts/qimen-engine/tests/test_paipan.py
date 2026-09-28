import importlib.util
import json
import subprocess
import unittest
from pathlib import Path


MODULE_PATH = Path(__file__).resolve().parents[1] / "paipan.py"
SPEC = importlib.util.spec_from_file_location("qimen_paipan", MODULE_PATH)
paipan = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(paipan)


class KinQiMenChartTest(unittest.TestCase):
    def test_reference_sample_hour_chart(self):
        local = paipan.parse_local_time("2020-10-07T14:23", "Asia/Shanghai")
        chart = paipan.build_chart(local, 3)
        self.assertEqual(chart["raw"]["排局"], "陰遁六局上元")
        self.assertEqual(chart["locked_palace"], 3)
        self.assertEqual(chart["palaces"]["3"]["spirit"], "六合")
        self.assertEqual(chart["palaces"]["3"]["door"], "死")
        self.assertEqual(chart["palaces"]["3"]["stem_pair"], "壬+辛")
        self.assertEqual(chart["method"], "时家奇门/阳盘/转盘/置闰")
        self.assertTrue(chart["palaces"]["7"]["day_void"])

    def test_reference_examples_show_method_disagreement(self):
        # Fixed comparison cases show a difference between the two rule engines.
        examples = (
            ("2020-05-02T19:00", "陽遁八局下元", "3", "休", "癸+壬"),
            ("2020-04-18T14:00", "陽遁七局下元", "7", "生", "乙+戊"),
        )
        for moment, expected_group, palace, door, pair in examples:
            with self.subTest(moment=moment):
                local = paipan.parse_local_time(moment, "Asia/Shanghai")
                chaibu = paipan.build_chart(local, None, "拆补")
                zhirun = paipan.build_chart(local, None, "置闰")
                self.assertEqual(chaibu["raw"]["排局"], expected_group)
                self.assertEqual(chaibu["palaces"][palace]["door"], door)
                self.assertEqual(chaibu["palaces"][palace]["stem_pair"], pair)
                self.assertNotEqual(zhirun["raw"]["排局"], expected_group)
                self.assertEqual(chaibu["engine_method_number"], 1)
                self.assertEqual(zhirun["engine_method_number"], 2)

    def test_engine_specific_spirit_names_are_preserved(self):
        local = paipan.parse_local_time("2020-05-02T19:00", "Asia/Shanghai")
        chart = paipan.build_chart(local, None, "拆补")
        self.assertEqual(chart["palaces"]["4"]["spirit"], "勾陈")
        self.assertEqual(chart["palaces"]["9"]["spirit"], "朱雀")
        self.assertEqual(chart["raw"]["神"]["巽"], "勾")

    def test_center_number_redirect_is_explicit(self):
        local = paipan.parse_local_time("2020-10-07T14:23", "Asia/Shanghai")
        chart = paipan.build_chart(local, 5)
        self.assertEqual(chart["reported_number"], 5)
        self.assertEqual(chart["locked_palace"], 2)
        self.assertTrue(chart["five_redirected_to_two"])


@unittest.skipUnless(
    (MODULE_PATH.parent / ".bin" / "atopx-qimen").is_file(),
    "先运行 scripts/qimen-engine/setup-atopx.sh",
)
class AtopxChartTest(unittest.TestCase):
    def test_three_reference_charts_use_the_labelled_zhirun_method(self):
        examples = (
            ("2020-10-07T14:23", "阴遁六局上元", "3", "死", "壬+辛", "六合"),
            ("2020-04-18T14:00", "阳遁七局下元", "7", "生", "乙+戊", "九天"),
            ("2020-05-02T19:00", "阳遁八局下元", "3", "休", "癸+壬", "六合"),
        )
        for moment, group, palace, door, pair, spirit in examples:
            with self.subTest(moment=moment):
                local = paipan.parse_local_time(moment, "Asia/Shanghai")
                chart = paipan.build_atopx_chart(local, None)
                self.assertEqual(chart["raw"]["排局"], group)
                self.assertEqual(chart["palaces"][palace]["door"], door)
                self.assertEqual(chart["palaces"][palace]["stem_pair"], pair)
                self.assertEqual(chart["palaces"][palace]["spirit"], spirit)
                self.assertEqual(chart["method"], "时家奇门/阳盘/转盘/置闰")

    def test_time_precision_is_not_chart_precision(self):
        first = paipan.build_atopx_chart(
            paipan.parse_local_time("2026-09-21T14:00:00", "Asia/Shanghai"), None
        )
        last = paipan.build_atopx_chart(
            paipan.parse_local_time("2026-09-21T14:59:59", "Asia/Shanghai"), None
        )
        next_hour = paipan.build_atopx_chart(
            paipan.parse_local_time("2026-09-21T15:00:00", "Asia/Shanghai"), None
        )
        self.assertEqual(first["palaces"], last["palaces"])
        self.assertNotEqual(first["palaces"], next_hour["palaces"])
        self.assertTrue(last["source_time"].startswith("2026-09-21T14:59:59"))

    def test_fixed_dates_change_the_dui_door(self):
        without = paipan.build_atopx_chart(
            paipan.parse_local_time("2020-04-09T19:12", "Asia/Shanghai"), None
        )
        with_rest = paipan.build_atopx_chart(
            paipan.parse_local_time("2020-04-13T19:12", "Asia/Shanghai"), None
        )
        self.assertNotEqual(without["palaces"]["7"]["door"], "休")
        self.assertEqual(with_rest["palaces"]["7"]["door"], "休")

    def test_fly_chart_keeps_center_and_missing_fields_explicit(self):
        local = paipan.parse_local_time("2020-10-07T14:23", "Asia/Shanghai")
        chart = paipan.build_atopx_chart(local, None, style="飞盘")
        self.assertEqual(sum(p["star"] is not None for p in chart["palaces"].values()), 9)
        self.assertEqual(sum(p["door"] is not None for p in chart["palaces"].values()), 8)
        self.assertIsNone(chart["locked_palace"])
        self.assertIn("尚未", chart["validation_scope"])
        with self.assertRaises(ValueError):
            paipan.build_atopx_chart(local, 3, style="飞盘")


@unittest.skipUnless(
    (MODULE_PATH.parent / ".ke-src" / "kinqimen.py").is_file()
    and (MODULE_PATH.parent / ".venv" / "bin" / "python").is_file(),
    "先运行 scripts/qimen-engine/setup-ke.sh",
)
class KeChartTest(unittest.TestCase):
    def test_ten_minute_boundary_changes_palace_chart(self):
        dates = ("2026-09-21", "2026-12-21", "2020-04-18")
        for day in dates:
            with self.subTest(day=day):
                before = paipan.build_ke_chart(paipan.parse_local_time(day + "T14:00:00", "Asia/Shanghai"))
                last = paipan.build_ke_chart(paipan.parse_local_time(day + "T14:09:59", "Asia/Shanghai"))
                next_ke = paipan.build_ke_chart(paipan.parse_local_time(day + "T14:10:00", "Asia/Shanghai"))
                self.assertEqual(before["palaces"], last["palaces"])
                self.assertNotEqual(before["palaces"], next_ke["palaces"])
                self.assertEqual(before["chart_granularity_minutes"], 10)
                self.assertEqual(before["ke_interval"]["end_exclusive"], next_ke["ke_interval"]["start_inclusive"])
                self.assertEqual(before["void_branches"]["minute"], before["raw"]["旬空"]["時空"])
                self.assertIsNone(before["void_branches"]["day"])
                self.assertIn("专用断法尚未核验", before["validation_scope"])

    def test_ke_is_explicit_and_cannot_accept_hour_number(self):
        python = MODULE_PATH.parent / ".venv" / "bin" / "python"
        base = [str(python), str(MODULE_PATH), "--family", "刻家", "--datetime", "2026-09-21T14:23"]
        result = subprocess.run(base, capture_output=True, text=True)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(json.loads(result.stdout)["method"], "刻家奇门/10分钟/转盘/上游标注置闰")
        self.assertEqual(json.loads(result.stdout)["route"]["engine"], "kinqimen-ke")
        rejected = subprocess.run(base + ["--number", "3"], capture_output=True, text=True)
        self.assertEqual(rejected.returncode, 2)
        self.assertIn("不套用", rejected.stderr)
        with self.assertRaises(ValueError):
            paipan.build_ke_chart(paipan.parse_local_time("2026-09-21T14:23", "Asia/Shanghai"), "拆补")
        with self.assertRaises(ValueError):
            paipan.build_ke_chart(paipan.parse_local_time("2026-09-21T14:23", "America/New_York"))

    def test_ke_dun_uses_a_different_route_from_hour_chart(self):
        python = MODULE_PATH.parent / ".venv" / "bin" / "python"
        base = [str(python), str(MODULE_PATH), "--datetime", "2026-12-21T14:20"]
        hour = json.loads(subprocess.check_output(base))
        ke = json.loads(subprocess.check_output(base + ["--family", "刻家"]))
        self.assertEqual(hour["route"]["yin_yang_dun"], "阳遁")
        self.assertEqual(ke["route"]["yin_yang_dun"], "阴遁")
        self.assertIn("时支", ke["route"]["yin_yang_source"])


class RoutingTest(unittest.TestCase):
    def test_default_cli_uses_reference_checked_hour_chart(self):
        result = subprocess.run(
            ["python3", str(MODULE_PATH), "--datetime", "2020-10-07T14:23"],
            capture_output=True, text=True,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        chart = json.loads(result.stdout)
        self.assertEqual(chart["engine"], "atopx/qimen")
        self.assertEqual(chart["route"]["family"], "时家")
        self.assertEqual(chart["route"]["engine"], "atopx")
        self.assertEqual(chart["route"]["yin_yang_dun"], "阴遁")

    def test_family_engine_and_method_boundaries(self):
        examples = (
            ("时家", "转盘", "置闰", None, "auto", "atopx"),
            ("时家", "转盘", "置闰", 5, "auto", "atopx"),
            ("时家", "飞盘", "置闰", None, "auto", "atopx"),
            ("刻家", "转盘", "置闰", None, "auto", "kinqimen-ke"),
            ("日家", "转盘", "拆补", None, "auto", "atopx"),
        )
        for family, style, method, number, engine, expected in examples:
            with self.subTest(family=family, style=style):
                self.assertEqual(paipan.resolve_route(family, style, method, number, engine), expected)
        invalid = (
            ("刻家", "转盘", "置闰", 5, "auto"),
            ("刻家", "转盘", "拆补", None, "auto"),
            ("刻家", "飞盘", "置闰", None, "auto"),
            ("时家", "转盘", "置闰", None, "kinqimen-ke"),
            ("时家", "飞盘", "置闰", None, "kinqimen"),
            ("日家", "转盘", "置闰", 5, "auto"),
            ("月家", "转盘", "拆补", None, "auto"),
        )
        for route in invalid:
            with self.subTest(route=route):
                with self.assertRaises(ValueError):
                    paipan.resolve_route(*route)


if __name__ == "__main__":
    unittest.main()
