import importlib.util
import pathlib
import unittest


SCRIPT = pathlib.Path(__file__).resolve().parents[1] / "phone_palace.py"
spec = importlib.util.spec_from_file_location("phone_palace", SCRIPT)
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class PhonePalaceTests(unittest.TestCase):
    def test_reference_example(self):
        actual = module.build_palace("0933377287")
        self.assertEqual(
            [actual[key] for key in ("palace", "spirit", "star", "door", "heaven_stem", "earth_stem")],
            ["震3", "九地", "天柱", "死门", "辛", "庚"],
        )

    def test_zero_is_position_specific(self):
        actual = module.build_palace("000000")
        self.assertEqual(actual["palace"], "宫空亡")
        self.assertEqual(actual["door"], "门空亡")
        self.assertEqual(actual["heaven_stem"], "癸")
        self.assertEqual(actual["earth_stem"], "癸")

    def test_requires_six_ascii_digits(self):
        for number in ("12345", "１２３４５６", "123-456"):
            with self.subTest(number=number), self.assertRaises(ValueError):
                module.build_palace(number)


if __name__ == "__main__":
    unittest.main()
