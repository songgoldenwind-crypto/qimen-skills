package almanac

import (
	"math"
	"strings"
)

// 寿星万年历 astronomical core (private).
//
// Implements solar-term and new-moon Julian Day calculations
// using truncated VSOP87 (xl0/xl1), IAU 1980 nutation (nutB),
// and a piecewise ΔT model (dtAt).  All identifiers are private:
// the package exposes only high-level types (Term, LunarDay, ...).

const (
	pi2          = 2 * math.Pi
	oneThird     = float64(1) / 3
	secondPerDay = 86400
	secondPerRad = 180 * 3600 / math.Pi
	// j2000 Julian Day of 2000-01-01 12:00 UT (epoch).
	j2000 = 2451545
)

// nutB IAU 1980 nutation series (10 leading terms × 5 floats).
var nutB = []float64{
	2.1824, -33.75705, 36e-6, -1720, 920,
	3.5069, 1256.66393, 11e-6, -132, 57,
	1.3375, 16799.4182, -51e-6, -23, 10,
	4.3649, -67.5141, 72e-6, 21, -9,
	0.04, -628.302, 0, -14, 0,
	2.36, 8328.691, 0, 7, 0,
	3.46, 1884.966, 0, -5, 2,
	5.44, 16833.175, 0, -4, 2,
	3.69, 25128.110, 0, -3, 0,
	3.55, 628.362, 0, 2, 0,
}

// dtAt piecewise ΔT polynomial coefficients
// (year, c0, c1, c2, c3) anchored every 100..300 years
// up to 1980, then every 4 years to 2048 + 2050 anchor.
// Predictions ≥ 2016 are fitted from skyfield's DE440s ΔT.
var dtAt = []float64{
	-4000, 108371.7, -13036.80, 392.000, 0.0000,
	-500, 17201.0, -627.82, 16.170, -0.3413,
	-150, 12200.6, -346.41, 5.403, -0.1593,
	150, 9113.8, -328.13, -1.647, 0.0377,
	500, 5707.5, -391.41, 0.915, 0.3145,
	900, 2203.4, -283.45, 13.034, -0.1778,
	1300, 490.1, -57.35, 2.085, -0.0072,
	1600, 120.0, -9.81, -1.532, 0.1403,
	1700, 10.2, -0.91, 0.510, -0.0370,
	1800, 13.4, -0.72, 0.202, -0.0193,
	1830, 7.8, -1.81, 0.416, -0.0247,
	1860, 8.3, -0.13, -0.406, 0.0292,
	1880, -5.4, 0.32, -0.183, 0.0173,
	1900, -2.3, 2.06, 0.169, -0.0135,
	1920, 21.2, 1.69, -0.304, 0.0167,
	1940, 24.2, 1.22, -0.064, 0.0031,
	1960, 33.2, 0.51, 0.231, -0.0109,
	1980, 51.0, 1.29, -0.026, 0.0032,
	2000, 63.87, 0.1, 0, 0,
	2005, 64.7, 0.21, 0, 0,
	2012, 66.8, 0.22, 0, 0,
	2016, 68.1024, 0.5456, -0.0542, -0.001172,
	2020, 69.3612, 0.0422, -0.0502, 0.006216,
	2024, 69.1752, -0.0335, -0.0048, 0.000811,
	2028, 69.0206, -0.0275, 0.0055, -0.000014,
	2032, 68.9981, 0.0163, 0.0054, 0.000006,
	2036, 69.1498, 0.0599, 0.0053, 0.000026,
	2040, 69.4751, 0.1035, 0.0051, 0.000046,
	2044, 69.9737, 0.1469, 0.0050, 0.000066,
	2048, 70.6451, 0.1903, 0.0049, 0.000085,
	2050, 71.0457,
}

// qiKb piecewise table of historical 节气 corrections (Julian Day, step).
// Used for the period 元嘉 (~AD 100) .. 1913 where high-precision VSOP
// theory disagrees with the calendar by 1 day occasionally.
var qiKb = []float64{
	1640650.479938, 15.21842500,
	1642476.703182, 15.21874996,
	1683430.515601, 15.218750011,
	1752157.640664, 15.218749978,
	1807675.003759, 15.218620279,
	1883627.765182, 15.218612292,
	1907369.128100, 15.218449176,
	1936603.140413, 15.218425000,
	1939145.524180, 15.218466998,
	1947180.798300, 15.218524844,
	1964362.041824, 15.218533526,
	1987372.340971, 15.218513908,
	1999653.819126, 15.218530782,
	2007445.469786, 15.218535181,
	2021324.917146, 15.218526248,
	2047257.232342, 15.218519654,
	2070282.898213, 15.218425000,
	2073204.872850, 15.218515221,
	2080144.500926, 15.218530782,
	2086703.688963, 15.218523776,
	2110033.182763, 15.218425000,
	2111190.300888, 15.218425000,
	2113731.271005, 15.218515671,
	2120670.840263, 15.218425000,
	2123973.309063, 15.218425000,
	2125068.997336, 15.218477932,
	2136026.312633, 15.218472436,
	2156099.495538, 15.218425000,
	2159021.324663, 15.218425000,
	2162308.575254, 15.218461742,
	2178485.706538, 15.218425000,
	2178759.662849, 15.218445786,
	2185334.020800, 15.218425000,
	2187525.481425, 15.218425000,
	2188621.191481, 15.218437494,
	2322147.76,
}

// shuoKb piecewise table for historical 朔望 corrections.
var shuoKb = []float64{
	1457698.231017, 29.53067166,
	1546082.512234, 29.53085106,
	1640640.735300, 29.53060000,
	1642472.151543, 29.53085439,
	1683430.509300, 29.53086148,
	1752148.041079, 29.53085097,
	1807665.420323, 29.53059851,
	1883618.114100, 29.53060000,
	1907360.704700, 29.53060000,
	1936596.224900, 29.53060000,
	1939135.675300, 29.53060000,
	1947168.00,
}

// qB / sB run-length-encoded historical ±1 day corrections
// for 节气 / 朔望 transitions across the entire span.
var qB = decode("FrcFs22AFsckF2tsDtFqEtF1posFdFgiFseFtmelpsEfhkF2anmelpFlF1ikrotcnEqEq2FfqmcDsrFor22FgFrcgDscFs22FgEeFtE2sfFs22sCoEsaF2tsD1FpeE2eFsssEciFsFnmelpFcFhkF2tcnEqEpFgkrotcnEqrEtFermcDsrE222FgBmcmr22DaEfnaF222sD1FpeForeF2tssEfiFpEoeFssD1iFstEqFppDgFstcnEqEpFg11FscnEqrAoAF2ClAEsDmDtCtBaDlAFbAEpAAAAAD2FgBiBqoBbnBaBoAAAAAAAEgDqAdBqAFrBaBoACdAAf1AACgAAAeBbCamDgEifAE2AABa1C1BgFdiAAACoCeE1ADiEifDaAEqAAFe1AcFbcAAAAAF1iFaAAACpACmFmAAAAAAAACrDaAAADG0")

var sB = decode("EqoFscDcrFpmEsF2DfFideFelFpFfFfFiaipqti1ksttikptikqckstekqttgkqttgkqteksttikptikq2fjstgjqttjkqttgkqtekstfkptikq2tijstgjiFkirFsAeACoFsiDaDiADc1AFbBfgdfikijFifegF1FhaikgFag1E2btaieeibggiffdeigFfqDfaiBkF1kEaikhkigeidhhdiegcFfakF1ggkidbiaedksaFffckekidhhdhdikcikiakicjF1deedFhFccgicdekgiFbiaikcfi1kbFibefgEgFdcFkFeFkdcfkF1kfkcickEiFkDacFiEfbiaejcFfffkhkdgkaiei1ehigikhdFikfckF1dhhdikcfgjikhfjicjicgiehdikcikggcifgiejF1jkieFhegikggcikFegiegkfjebhigikggcikdgkaFkijcfkcikfkcifikiggkaeeigefkcdfcfkhkdgkegieidhijcFfakhfgeidieidiegikhfkfckfcjbdehdikggikgkfkicjicjF1dbidikFiggcifgiejkiegkigcdiegfggcikdbgfgefjF1kfegikggcikdgFkeeijcfkcikfkekcikdgkabhkFikaffcfkhkdgkegbiaekfkiakicjhfgqdq2fkiakgkfkhfkfcjiekgFebicggbedF1jikejbbbiakgbgkacgiejkijjgigfiakggfggcibFifjefjF1kfekdgjcibFeFkijcfkfhkfkeaieigekgbhkfikidfcjeaibgekgdkiffiffkiakF1jhbakgdki1dj1ikfkicjicjieeFkgdkicggkighdF1jfgkgfgbdkicggfggkidFkiekgijkeigfiskiggfaidheigF1jekijcikickiggkidhhdbgcfkFikikhkigeidieFikggikhkffaffijhidhhakgdkhkijF1kiakF1kfheakgdkifiggkigicjiejkieedikgdfcggkigieeiejfgkgkigbgikicggkiaideeijkefjeijikhkiggkiaidheigcikaikffikijgkiahi1hhdikgjfifaakekighie1hiaikggikhkffakicjhiahaikggikhkijF1kfejfeFhidikggiffiggkigicjiekgieeigikggiffiggkidheigkgfjkeigiegikifiggkidhedeijcfkFikikhkiggkidhh1ehigcikaffkhkiggkidhh1hhigikekfiFkFikcidhh1hitcikggikhkfkicjicghiediaikggikhkijbjfejfeFhaikggifikiggkigiejkikgkgieeigikggiffiggkigieeigekijcijikggifikiggkideedeijkefkfckikhkiggkidhh1ehijcikaffkhkiggkidhh1hhigikhkikFikfckcidhh1hiaikgjikhfjicjicgiehdikcikggifikigiejfejkieFhegikggifikiggfghigkfjeijkhigikggifikiggkigieeijcijcikfksikifikiggkidehdeijcfdckikhkiggkhghh1ehijikifffffkhsFngErD1pAfBoDd1BlEtFqA2AqoEpDqElAEsEeB2BmADlDkqBtC1FnEpDqnEmFsFsAFnllBbFmDsDiCtDmAB2BmtCgpEplCpAEiBiEoFqFtEqsDcCnFtADnFlEgdkEgmEtEsCtDmADqFtAFrAtEcCqAE1BoFqC1F1DrFtBmFtAC2ACnFaoCgADcADcCcFfoFtDlAFgmFqBq2bpEoAEmkqnEeCtAE1bAEqgDfFfCrgEcBrACfAAABqAAB1AAClEnFeCtCgAADqDoBmtAAACbFiAAADsEtBqAB2FsDqpFqEmFsCeDtFlCeDtoEpClEqAAFrAFoCgFmFsFqEnAEcCqFeCtFtEnAEeFtAAEkFnErAABbFkADnAAeCtFeAfBoAEpFtAABtFqAApDcCGJ")

// decode expands the run-length-encoded corrections strings into a 0/1/2 digit sequence.
func decode(s string) string {
	o := "0000000000"
	o2 := o + o
	s = strings.Replace(s, "J", "00", -1)
	s = strings.Replace(s, "I", "000", -1)
	s = strings.Replace(s, "H", "0000", -1)
	s = strings.Replace(s, "G", "00000", -1)
	s = strings.Replace(s, "t", "02", -1)
	s = strings.Replace(s, "s", "002", -1)
	s = strings.Replace(s, "r", "0002", -1)
	s = strings.Replace(s, "q", "00002", -1)
	s = strings.Replace(s, "p", "000002", -1)
	s = strings.Replace(s, "o", "0000002", -1)
	s = strings.Replace(s, "n", "00000002", -1)
	s = strings.Replace(s, "m", "000000002", -1)
	s = strings.Replace(s, "l", "0000000002", -1)
	s = strings.Replace(s, "k", "01", -1)
	s = strings.Replace(s, "j", "0101", -1)
	s = strings.Replace(s, "i", "001", -1)
	s = strings.Replace(s, "h", "001001", -1)
	s = strings.Replace(s, "g", "0001", -1)
	s = strings.Replace(s, "f", "00001", -1)
	s = strings.Replace(s, "e", "000001", -1)
	s = strings.Replace(s, "d", "0000001", -1)
	s = strings.Replace(s, "c", "00000001", -1)
	s = strings.Replace(s, "b", "000000001", -1)
	s = strings.Replace(s, "a", "0000000001", -1)
	s = strings.Replace(s, "A", o2+o2+o2, -1)
	s = strings.Replace(s, "B", o2+o2+o, -1)
	s = strings.Replace(s, "C", o2+o2, -1)
	s = strings.Replace(s, "D", o2+o, -1)
	s = strings.Replace(s, "E", o2, -1)
	s = strings.Replace(s, "F", o, -1)
	return s
}

// nutationLon2 truncated IAU 1980 nutation longitude (radians).
func nutationLon2(t float64) float64 {
	a := -1.742 * t
	t2 := t * t
	dl := float64(0)
	j := len(nutB)
	for i := 0; i < j; i += 5 {
		dl += (nutB[i+3] + a) * math.Sin(nutB[i]+nutB[i+1]*t+nutB[i+2]*t2)
		a = 0
	}
	return dl / 100 / secondPerRad
}

// eLon truncated VSOP87 ecliptic longitude of Earth (radians).
// n controls truncation depth: <0 = full, otherwise proportional.
func eLon(t float64, n int) float64 {
	t /= 10
	v := float64(0)
	tn := float64(1)
	pn := 1
	m0 := xl0[pn+1] - xl0[pn]
	for i := 0; i < 6; i++ {
		n1 := int(xl0[pn+i])
		n2 := int(xl0[pn+1+i])
		n0 := n2 - n1
		if n0 == 0 {
			continue
		}
		m := 0
		if n < 0 {
			m = n2
		} else {
			m = int(float64(3*n*n0)/m0+0.5) + n1
			if i != 0 {
				m += 3
			}
			if m > n2 {
				m = n2
			}
		}
		c := float64(0)
		for j := n1; j < m; j += 3 {
			c += xl0[j] * math.Cos(xl0[j+1]+t*xl0[j+2])
		}
		v += c * tn
		tn = tn * t
	}
	v /= xl0[0]
	t2 := t * t
	t3 := t2 * t
	return v + (-0.0728-2.7702*t-1.1019*t2-0.0996*t3)/secondPerRad
}

// mLon Moon ecliptic longitude (radians) — uses xl1 table.
func mLon(t float64, n int) float64 {
	obl := len(xl1[0])
	tn := float64(1)
	v := float64(0)
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t
	tx := t - 10
	v += (3.81034409 + 8399.684730072*t - 3.319e-05*t2 + 3.11e-08*t3 - 2.033e-10*t4) * secondPerRad
	v += 5028.792262*t + 1.1124406*t2 + 0.00007699*t3 - 0.000023479*t4 - 0.0000000178*t5
	if tx > 0 {
		v += -0.866 + 1.43*tx + 0.054*tx*tx
	}
	t2 /= 1e4
	t3 /= 1e8
	t4 /= 1e8

	n *= 6
	if n < 0 {
		n = obl
	}
	x := len(xl1)
	for i := 0; i < x; i++ {
		f := xl1[i][:]
		l := len(f)
		m := int(float64(n*l)/float64(obl) + 0.5)
		if i > 0 {
			m += 6
		}
		if m >= l {
			m = l
		}
		c := float64(0)
		for j := 0; j < m; j += 6 {
			c += f[j] * math.Cos(f[j+1]+t*f[j+2]+t2*f[j+3]+t3*f[j+4]+t4*f[j+5])
		}
		v += c * tn
		tn *= t
	}
	return v / secondPerRad
}

// gxcSunLon stellar aberration of the Sun longitude (radians).
func gxcSunLon(t float64) float64 {
	t2 := t * t
	return -20.49552 * (1 + (0.016708634-0.000042037*t-0.0000001267*t2)*math.Cos(-0.043126+628.301955*t-0.000002732*t2)) / secondPerRad
}

// ev Earth orbital angular velocity correction.
func ev(t float64) float64 {
	f := 628.307585 * t
	return 628.332 + 21*math.Sin(1.527+f) + 0.44*math.Sin(1.48+f*2) + 0.129*math.Sin(5.82+f)*t + 0.00055*math.Sin(4.21+f)*t*t
}

// saLon apparent geocentric Sun longitude (radians).
func saLon(t float64, n int) float64 {
	return eLon(t, n) + nutationLon2(t) + gxcSunLon(t) + math.Pi
}

// dtExt long-tail ΔT extrapolation (jsd = squared annual drift).
func dtExt(y, jsd float64) float64 {
	dy := (y - 1820) / 100
	return -20 + jsd*dy*dy
}

// dtCalc ΔT in seconds for a given decimal year (gregorian).
func dtCalc(y float64) float64 {
	size := len(dtAt)
	y0 := dtAt[size-2]
	t0 := dtAt[size-1]
	if y >= y0 {
		jsd := float64(31)
		if y > y0+100 {
			return dtExt(y, jsd)
		}
		v := dtExt(y, jsd)
		dv := dtExt(y0, jsd) - t0
		return v - dv*(y0+100-y)/100
	}
	i := 0
	for ; i < size; i += 5 {
		if y < dtAt[i+5] {
			break
		}
	}
	t1 := (y - dtAt[i]) / (dtAt[i+5] - dtAt[i]) * 10
	t2 := t1 * t1
	t3 := t2 * t1
	return dtAt[i+1] + dtAt[i+2]*t1 + dtAt[i+3]*t2 + dtAt[i+4]*t3
}

// dtT ΔT in days for a Julian century offset from J2000.
func dtT(t float64) float64 {
	return dtCalc(t/365.2425+2000) / secondPerDay
}

// mv Moon orbital angular velocity correction.
func mv(t float64) float64 {
	v := 8399.71 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t*t)
	return v - (179*math.Sin(2.543+15542.7543*t) + 160*math.Sin(0.1874+7214.0629*t) + 62*math.Sin(3.14+16657.3828*t) + 34*math.Sin(4.827+16866.9323*t) + 22*math.Sin(4.9+23871.4457*t) + 12*math.Sin(2.59+14914.4523*t) + 7*math.Sin(0.23+6585.7609*t) + 5*math.Sin(0.9+25195.624*t) + 5*math.Sin(2.32-7700.3895*t) + 5*math.Sin(3.88+8956.9934*t) + 5*math.Sin(0.49+7771.3771*t))
}

// saLonT inverse: longitude → Julian century for the Sun.
func saLonT(w float64) float64 {
	v := 628.3319653318
	t := (w - 1.75347 - math.Pi) / v
	v = ev(t)
	t += (w - saLon(t, 10)) / v
	return t + (w-saLon(t, -1))/ev(t)
}

// msaLon Moon - Sun apparent longitude difference (radians).
func msaLon(t float64, mn, sn int) float64 {
	return mLon(t, mn) + (-3.4e-6) - (eLon(t, sn) + gxcSunLon(t) + math.Pi)
}

// msaLonT inverse: lunation phase → Julian century.
func msaLonT(w float64) float64 {
	v := 7771.37714500204
	t := (w + 1.08472) / v
	t += (w - msaLon(t, 3, 3)) / v
	v = mv(t) - ev(t)
	t += (w - msaLon(t, 20, 10)) / v
	return t + (w-msaLon(t, -1, 60))/v
}

// saLonT2 low-precision longitude → Julian century for the Sun.
func saLonT2(w float64) float64 {
	v := 628.3319653318
	t := (w - 1.75347 - math.Pi) / v
	t -= (0.000005297*t*t + 0.0334166*math.Cos(4.669257+628.307585*t) + 0.0002061*math.Cos(2.67823+628.307585*t)*t) / v
	return t + (w-eLon(t, 8)-math.Pi+(20.5+17.2*math.Sin(2.1824-33.75705*t))/secondPerRad)/v
}

// msaLonT2 low-precision lunation → Julian century.
func msaLonT2(w float64) float64 {
	v := 7771.37714500204
	t := (w + 1.08472) / v
	t2 := t * t
	t -= (-0.00003309*t2 + 0.10976*math.Cos(0.784758+8328.6914246*t+0.000152292*t2) + 0.02224*math.Cos(0.18740+7214.0628654*t-0.00021848*t2) - 0.03342*math.Cos(4.669257+628.307585*t)) / v
	t2 = t * t
	l := mLon(t, 20) - (4.8950632 + 628.3319653318*t + 0.000005297*t2 + 0.0334166*math.Cos(4.669257+628.307585*t) + 0.0002061*math.Cos(2.67823+628.307585*t)*t + 0.000349*math.Cos(4.6261+1256.61517*t) - 20.5/secondPerRad)
	v = 7771.38 - 914*math.Sin(0.7848+8328.691425*t+0.0001523*t2) - 179*math.Sin(2.543+15542.7543*t) - 160*math.Sin(0.1874+7214.0629*t)
	return t + (w-l)/v
}

// qiHigh high-precision 节气 Julian Day (offset from J2000).
func qiHigh(w float64) float64 {
	t := saLonT2(w) * 36525
	t = t - dtT(t) + oneThird
	v := math.Mod(t+0.5, 1) * secondPerDay
	if v < 1200 || v > secondPerDay-1200 {
		t = saLonT(w)*36525 - dtT(t) + oneThird
	}
	return t
}

// shuoHigh high-precision 朔 Julian Day (offset from J2000).
func shuoHigh(w float64) float64 {
	t := msaLonT2(w) * 36525
	t = t - dtT(t) + oneThird
	v := math.Mod(t+0.5, 1) * secondPerDay
	if v < 1800 || v > secondPerDay-1800 {
		t = msaLonT(w)*36525 - dtT(t) + oneThird
	}
	return t
}

// qiLow low-precision 节气 fallback for very old / very far dates.
func qiLow(w float64) float64 {
	v := 628.3319653318
	t := (w - 4.895062166) / v
	t -= (53*t*t + 334116*math.Cos(4.67+628.307585*t) + 2061*math.Cos(2.678+628.3076*t)*t) / v / 10000000
	n := 48950621.66 + 6283319653.318*t + 53*t*t + 334166*math.Cos(4.669257+628.307585*t) + 3489*math.Cos(4.6261+1256.61517*t) + 2060.6*math.Cos(2.67823+628.307585*t)*t - 994 - 834*math.Sin(2.1824-33.75705*t)
	t -= (n/10000000-w)/628.332 + (32*(t+1.8)*(t+1.8)-20)/secondPerDay/36525
	return t*36525 + oneThird
}

// shuoLow low-precision 朔 fallback.
func shuoLow(w float64) float64 {
	v := 7771.37714500204
	t := (w + 1.08472) / v
	t -= (-0.0000331*t*t+0.10976*math.Cos(0.785+8328.6914*t)+0.02224*math.Cos(0.187+7214.0629*t)-0.03342*math.Cos(4.669+628.3076*t))/v + (32*(t+1.8)*(t+1.8)-20)/secondPerDay/36525
	return t*36525 + oneThird
}

// calcShuo Julian-Day offset (from J2000) of the 朔 (new moon) immediately ≤ jd.
func calcShuo(jd float64) float64 {
	size := len(shuoKb)
	d := float64(0)
	i := 0
	pc := float64(14)
	jd += j2000
	f1 := shuoKb[0] - pc
	f2 := shuoKb[size-1] - pc
	f3 := float64(2436935)
	switch {
	case jd < f1 || jd >= f3:
		d = math.Floor(shuoHigh(math.Floor((jd+pc-2451551)/29.5306)*pi2) + 0.5)
	case jd >= f1 && jd < f2:
		for i = 0; i < size; i += 2 {
			if jd+pc < shuoKb[i+2] {
				break
			}
		}
		d = shuoKb[i] + shuoKb[i+1]*math.Floor((jd+pc-shuoKb[i])/shuoKb[i+1])
		d = math.Floor(d + 0.5)
		if d == 1683460 {
			d++
		}
		d -= j2000
	case jd >= f2 && jd < f3:
		d = math.Floor(shuoLow(math.Floor((jd+pc-2451551)/29.5306)*pi2) + 0.5)
		from := int((jd - f2) / 29.5306)
		n := sB[from : from+1]
		switch n {
		case "1":
			d++
		case "2":
			d--
		}
	}
	return d
}

// calcQi Julian-Day offset (from J2000) of the 节气 (solar term) immediately ≤ jd.
func calcQi(jd float64) float64 {
	size := len(qiKb)
	d := float64(0)
	i := 0
	pc := float64(7)
	jd += j2000
	f1 := qiKb[0] - pc
	f2 := qiKb[size-1] - pc
	f3 := float64(2436935)
	switch {
	case jd < f1 || jd >= f3:
		d = math.Floor(qiHigh(math.Floor((jd+pc-2451259)/365.2422*24)*math.Pi/12) + 0.5)
	case jd >= f1 && jd < f2:
		for i = 0; i < size; i += 2 {
			if jd+pc < qiKb[i+2] {
				break
			}
		}
		d = qiKb[i] + qiKb[i+1]*math.Floor((jd+pc-qiKb[i])/qiKb[i+1])
		d = math.Floor(d + 0.5)
		if d == 1683460 {
			d++
		}
		d -= j2000
	case jd >= f2 && jd < f3:
		d = math.Floor(qiLow(math.Floor((jd+pc-2451259)/365.2422*24)*math.Pi/12) + 0.5)
		from := int((jd - f2) / 365.2422 * 24)
		n := qB[from : from+1]
		switch n {
		case "1":
			d++
		case "2":
			d--
		}
	}
	return d
}

// qiAccurate refines a 节气 result to second-level precision.
func qiAccurate(w float64) float64 {
	t := saLonT(w) * 36525
	return t - dtT(t) + oneThird
}

// qiAccurate2 refines based on a coarse Julian Day estimate.
func qiAccurate2(jd float64) float64 {
	d := math.Pi / 12
	w := math.Floor((jd+293)/365.2422*24) * d
	a := qiAccurate(w)
	if a-jd > 5 {
		return qiAccurate(w - d)
	}
	if a-jd < -5 {
		return qiAccurate(w + d)
	}
	return a
}
