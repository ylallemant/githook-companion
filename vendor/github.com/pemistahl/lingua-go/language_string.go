// Code generated by "stringer -type=Language"; DO NOT EDIT.

package lingua

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Afrikaans-0]
	_ = x[Albanian-1]
	_ = x[Arabic-2]
	_ = x[Armenian-3]
	_ = x[Azerbaijani-4]
	_ = x[Basque-5]
	_ = x[Belarusian-6]
	_ = x[Bengali-7]
	_ = x[Bokmal-8]
	_ = x[Bosnian-9]
	_ = x[Bulgarian-10]
	_ = x[Catalan-11]
	_ = x[Chinese-12]
	_ = x[Croatian-13]
	_ = x[Czech-14]
	_ = x[Danish-15]
	_ = x[Dutch-16]
	_ = x[English-17]
	_ = x[Esperanto-18]
	_ = x[Estonian-19]
	_ = x[Finnish-20]
	_ = x[French-21]
	_ = x[Ganda-22]
	_ = x[Georgian-23]
	_ = x[German-24]
	_ = x[Greek-25]
	_ = x[Gujarati-26]
	_ = x[Hebrew-27]
	_ = x[Hindi-28]
	_ = x[Hungarian-29]
	_ = x[Icelandic-30]
	_ = x[Indonesian-31]
	_ = x[Irish-32]
	_ = x[Italian-33]
	_ = x[Japanese-34]
	_ = x[Kazakh-35]
	_ = x[Korean-36]
	_ = x[Latin-37]
	_ = x[Latvian-38]
	_ = x[Lithuanian-39]
	_ = x[Macedonian-40]
	_ = x[Malay-41]
	_ = x[Maori-42]
	_ = x[Marathi-43]
	_ = x[Mongolian-44]
	_ = x[Nynorsk-45]
	_ = x[Persian-46]
	_ = x[Polish-47]
	_ = x[Portuguese-48]
	_ = x[Punjabi-49]
	_ = x[Romanian-50]
	_ = x[Russian-51]
	_ = x[Serbian-52]
	_ = x[Shona-53]
	_ = x[Slovak-54]
	_ = x[Slovene-55]
	_ = x[Somali-56]
	_ = x[Sotho-57]
	_ = x[Spanish-58]
	_ = x[Swahili-59]
	_ = x[Swedish-60]
	_ = x[Tagalog-61]
	_ = x[Tamil-62]
	_ = x[Telugu-63]
	_ = x[Thai-64]
	_ = x[Tsonga-65]
	_ = x[Tswana-66]
	_ = x[Turkish-67]
	_ = x[Ukrainian-68]
	_ = x[Urdu-69]
	_ = x[Vietnamese-70]
	_ = x[Welsh-71]
	_ = x[Xhosa-72]
	_ = x[Yoruba-73]
	_ = x[Zulu-74]
	_ = x[Unknown-75]
}

const _Language_name = "AfrikaansAlbanianArabicArmenianAzerbaijaniBasqueBelarusianBengaliBokmalBosnianBulgarianCatalanChineseCroatianCzechDanishDutchEnglishEsperantoEstonianFinnishFrenchGandaGeorgianGermanGreekGujaratiHebrewHindiHungarianIcelandicIndonesianIrishItalianJapaneseKazakhKoreanLatinLatvianLithuanianMacedonianMalayMaoriMarathiMongolianNynorskPersianPolishPortuguesePunjabiRomanianRussianSerbianShonaSlovakSloveneSomaliSothoSpanishSwahiliSwedishTagalogTamilTeluguThaiTsongaTswanaTurkishUkrainianUrduVietnameseWelshXhosaYorubaZuluUnknown"

var _Language_index = [...]uint16{0, 9, 17, 23, 31, 42, 48, 58, 65, 71, 78, 87, 94, 101, 109, 114, 120, 125, 132, 141, 149, 156, 162, 167, 175, 181, 186, 194, 200, 205, 214, 223, 233, 238, 245, 253, 259, 265, 270, 277, 287, 297, 302, 307, 314, 323, 330, 337, 343, 353, 360, 368, 375, 382, 387, 393, 400, 406, 411, 418, 425, 432, 439, 444, 450, 454, 460, 466, 473, 482, 486, 496, 501, 506, 512, 516, 523}

func (language Language) String() string {
	if language < 0 || language >= Language(len(_Language_index)-1) {
		return "Language(" + strconv.FormatInt(int64(language), 10) + ")"
	}
	return _Language_name[_Language_index[language]:_Language_index[language+1]]
}
