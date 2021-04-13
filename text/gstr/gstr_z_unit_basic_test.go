// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/basicfu/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/basicfu/gf/frame/g"
	"github.com/basicfu/gf/test/gtest"
)

func Test_Replace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(Replace(s1, "ab", "AB"), "ABcdEFG乱入的中文ABcdefg")
		t.Assert(Replace(s1, "EF", "ef"), "abcdefG乱入的中文abcdefg")
		t.Assert(Replace(s1, "MN", "mn"), s1)

		t.Assert(ReplaceByArray(s1, g.ArrayStr{
			"a", "A",
			"A", "-",
			"a",
		}), "-bcdEFG乱入的中文-bcdefg")

		t.Assert(ReplaceByMap(s1, g.MapStrStr{
			"a": "A",
			"G": "g",
		}), "AbcdEFg乱入的中文Abcdefg")
	})
}

func Test_ReplaceI_1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcd乱入的中文ABCD"
		s2 := "a"
		t.Assert(ReplaceI(s1, "ab", "aa"), "aacd乱入的中文aaCD")
		t.Assert(ReplaceI(s1, "ab", "aa", 0), "abcd乱入的中文ABCD")
		t.Assert(ReplaceI(s1, "ab", "aa", 1), "aacd乱入的中文ABCD")

		t.Assert(ReplaceI(s1, "abcd", "-"), "-乱入的中文-")
		t.Assert(ReplaceI(s1, "abcd", "-", 1), "-乱入的中文ABCD")

		t.Assert(ReplaceI(s1, "abcd乱入的", ""), "中文ABCD")
		t.Assert(ReplaceI(s1, "ABCD乱入的", ""), "中文ABCD")

		t.Assert(ReplaceI(s2, "A", "-"), "-")
		t.Assert(ReplaceI(s2, "a", "-"), "-")

		t.Assert(ReplaceIByArray(s1, g.ArrayStr{
			"abcd乱入的", "-",
			"-", "=",
			"a",
		}), "=中文ABCD")

		t.Assert(ReplaceIByMap(s1, g.MapStrStr{
			"ab": "-",
			"CD": "=",
		}), "-=乱入的中文-=")
	})
}

func Test_ToLower(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "abcdefg乱入的中文abcdefg"
		t.Assert(ToLower(s1), e1)
	})
}

func Test_ToUpper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "ABCDEFG乱入的中文ABCDEFG"
		t.Assert(ToUpper(s1), e1)
	})
}

func Test_UcFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "abcdEFG乱入的中文abcdefg"
		e1 := "AbcdEFG乱入的中文abcdefg"
		t.Assert(UcFirst(""), "")
		t.Assert(UcFirst(s1), e1)
		t.Assert(UcFirst(e1), e1)
	})
}

func Test_LcFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "AbcdEFG乱入的中文abcdefg"
		e1 := "abcdEFG乱入的中文abcdefg"
		t.Assert(LcFirst(""), "")
		t.Assert(LcFirst(s1), e1)
		t.Assert(LcFirst(e1), e1)
	})
}

func Test_UcWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := "我爱GF: i love go frame"
		e1 := "我爱GF: I Love Go Frame"
		t.Assert(UcWords(s1), e1)
	})
}

func Test_IsLetterLower(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsLetterLower('a'), true)
		t.Assert(IsLetterLower('A'), false)
		t.Assert(IsLetterLower('1'), false)
	})
}

func Test_IsLetterUpper(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsLetterUpper('a'), false)
		t.Assert(IsLetterUpper('A'), true)
		t.Assert(IsLetterUpper('1'), false)
	})
}

func Test_IsNumeric(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsNumeric("1a我"), false)
		t.Assert(IsNumeric("0123"), true)
		t.Assert(IsNumeric("我是中国人"), false)
	})
}

func Test_SubStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(SubStr("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(SubStr("我爱GoFrame", 6), "GoFrame")
		t.Assert(SubStr("我爱GoFrame", 6, 2), "Go")
		t.Assert(SubStr("我爱GoFrame", -1, 30), "我爱GoFrame")
		t.Assert(SubStr("我爱GoFrame", 30, 30), "")
	})
}

func Test_SubStrRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(SubStrRune("我爱GoFrame", 0), "我爱GoFrame")
		t.Assert(SubStrRune("我爱GoFrame", 2), "GoFrame")
		t.Assert(SubStrRune("我爱GoFrame", 2, 2), "Go")
		t.Assert(SubStrRune("我爱GoFrame", -1, 30), "我爱GoFrame")
		t.Assert(SubStrRune("我爱GoFrame", 30, 30), "")
	})
}

func Test_StrLimit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(StrLimit("我爱GoFrame", 6), "我爱...")
		t.Assert(StrLimit("我爱GoFrame", 6, ""), "我爱")
		t.Assert(StrLimit("我爱GoFrame", 6, "**"), "我爱**")
		t.Assert(StrLimit("我爱GoFrame", 8, ""), "我爱Go")
		t.Assert(StrLimit("*", 4, ""), "*")
	})
}

func Test_StrLimitRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(StrLimitRune("我爱GoFrame", 2), "我爱...")
		t.Assert(StrLimitRune("我爱GoFrame", 2, ""), "我爱")
		t.Assert(StrLimitRune("我爱GoFrame", 2, "**"), "我爱**")
		t.Assert(StrLimitRune("我爱GoFrame", 4, ""), "我爱Go")
		t.Assert(StrLimitRune("*", 4, ""), "*")
	})
}

func Test_HasPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(HasPrefix("我爱GoFrame", "我爱"), true)
		t.Assert(HasPrefix("en我爱GoFrame", "我爱"), false)
		t.Assert(HasPrefix("en我爱GoFrame", "en"), true)
	})
}

func Test_HasSuffix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(HasSuffix("我爱GoFrame", "GoFrame"), true)
		t.Assert(HasSuffix("en我爱GoFrame", "a"), false)
		t.Assert(HasSuffix("GoFrame很棒", "棒"), true)
	})
}

func Test_Reverse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Reverse("我爱123"), "321爱我")
	})
}

func Test_NumberFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(NumberFormat(1234567.8910, 2, ".", ","), "1,234,567.89")
		t.Assert(NumberFormat(1234567.8910, 2, "#", "/"), "1/234/567#89")
		t.Assert(NumberFormat(-1234567.8910, 2, "#", "/"), "-1/234/567#89")
	})
}

func Test_ChunkSplit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(ChunkSplit("1234", 1, "#"), "1#2#3#4#")
		t.Assert(ChunkSplit("我爱123", 1, "#"), "我#爱#1#2#3#")
		t.Assert(ChunkSplit("1234", 1, ""), "1\r\n2\r\n3\r\n4\r\n")
	})
}

func Test_SplitAndTrim(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `

010    

020  

`
		a := SplitAndTrim(s, "\n", "0")
		t.Assert(len(a), 2)
		t.Assert(a[0], "1")
		t.Assert(a[1], "2")
	})
}

func Test_Fields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Fields("我爱 Go Frame"), []string{
			"我爱", "Go", "Frame",
		})
	})
}

func Test_CountWords(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(CountWords("我爱 Go Go Go"), map[string]int{
			"Go": 3,
			"我爱": 1,
		})
	})
}

func Test_CountChars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(CountChars("我爱 Go Go Go"), map[string]int{
			" ": 3,
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
		t.Assert(CountChars("我爱 Go Go Go", true), map[string]int{
			"G": 3,
			"o": 3,
			"我": 1,
			"爱": 1,
		})
	})
}

func Test_WordWrap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(WordWrap("12 34", 2, "<br>"), "12<br>34")
		t.Assert(WordWrap("12 34", 2, "\n"), "12\n34")
		t.Assert(WordWrap("我爱 GF", 2, "\n"), "我爱\nGF")
		t.Assert(WordWrap("A very long woooooooooooooooooord. and something", 7, "<br>"),
			"A very<br>long<br>woooooooooooooooooord.<br>and<br>something")
	})
}

func Test_RuneLen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(RuneLen("1234"), 4)
		t.Assert(RuneLen("我爱GoFrame"), 9)
	})
}

func Test_Repeat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Repeat("go", 3), "gogogo")
		t.Assert(Repeat("好的", 3), "好的好的好的")
	})
}

func Test_Str(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Str("name@example.com", "@"), "@example.com")
		t.Assert(Str("name@example.com", ""), "")
		t.Assert(Str("name@example.com", "z"), "")
	})
}

func Test_Shuffle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(len(Shuffle("123456")), 6)
	})
}

func Test_Split(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Split("1.2", "."), []string{"1", "2"})
		t.Assert(Split("我爱 - GoFrame", " - "), []string{"我爱", "GoFrame"})
	})
}

func Test_Join(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Join([]string{"我爱", "GoFrame"}, " - "), "我爱 - GoFrame")
	})
}

func Test_Explode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Explode(" - ", "我爱 - GoFrame"), []string{"我爱", "GoFrame"})
	})
}

func Test_Implode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Implode(" - ", []string{"我爱", "GoFrame"}), "我爱 - GoFrame")
	})
}

func Test_Chr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Chr(65), "A")
	})
}

func Test_Ord(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Ord("A"), 65)
	})
}

func Test_HideStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(HideStr("15928008611", 40, "*"), "159****8611")
		t.Assert(HideStr("john@kohg.cn", 40, "*"), "jo*n@kohg.cn")
		t.Assert(HideStr("张三", 50, "*"), "张*")
		t.Assert(HideStr("张小三", 50, "*"), "张*三")
		t.Assert(HideStr("欧阳小三", 50, "*"), "欧**三")
	})
}

func Test_Nl2Br(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Nl2Br("1\n2"), "1<br>2")
		t.Assert(Nl2Br("1\r\n2"), "1<br>2")
		t.Assert(Nl2Br("1\r\n2", true), "1<br />2")
	})
}

func Test_AddSlashes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(AddSlashes(`1'2"3\`), `1\'2\"3\\`)
	})
}

func Test_StripSlashes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(StripSlashes(`1\'2\"3\\`), `1'2"3\`)
	})
}

func Test_QuoteMeta(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(QuoteMeta(`.\+*?[^]($)`), `\.\\\+\*\?\[\^\]\(\$\)`)
		t.Assert(QuoteMeta(`.\+*中国?[^]($)`), `\.\\\+\*中国\?\[\^\]\(\$\)`)
		t.Assert(QuoteMeta(`.''`, `'`), `.\'\'`)
		t.Assert(QuoteMeta(`中国.''`, `'`), `中国.\'\'`)
	})
}

func Test_Count(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "abcdaAD"
		t.Assert(Count(s, "0"), 0)
		t.Assert(Count(s, "a"), 2)
		t.Assert(Count(s, "b"), 1)
		t.Assert(Count(s, "d"), 1)
	})
}

func Test_CountI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "abcdaAD"
		t.Assert(CountI(s, "0"), 0)
		t.Assert(CountI(s, "a"), 3)
		t.Assert(CountI(s, "b"), 1)
		t.Assert(CountI(s, "d"), 2)
	})
}

func Test_Compare(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Compare("a", "b"), -1)
		t.Assert(Compare("a", "a"), 0)
		t.Assert(Compare("b", "a"), 1)
	})
}

func Test_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Equal("a", "A"), true)
		t.Assert(Equal("a", "a"), true)
		t.Assert(Equal("b", "a"), false)
	})
}

func Test_Contains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(Contains("abc", "a"), true)
		t.Assert(Contains("abc", "A"), false)
		t.Assert(Contains("abc", "ab"), true)
		t.Assert(Contains("abc", "abc"), true)
	})
}

func Test_ContainsI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(ContainsI("abc", "a"), true)
		t.Assert(ContainsI("abc", "A"), true)
		t.Assert(ContainsI("abc", "Ab"), true)
		t.Assert(ContainsI("abc", "ABC"), true)
		t.Assert(ContainsI("abc", "ABCD"), false)
		t.Assert(ContainsI("abc", "D"), false)
	})
}

func Test_ContainsAny(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(ContainsAny("abc", "a"), true)
		t.Assert(ContainsAny("abc", "cd"), true)
		t.Assert(ContainsAny("abc", "de"), false)
		t.Assert(ContainsAny("abc", "A"), false)
	})
}

func Test_SearchArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(SearchArray(a, "a"), 0)
		t.AssertEQ(SearchArray(a, "b"), 1)
		t.AssertEQ(SearchArray(a, "c"), 2)
		t.AssertEQ(SearchArray(a, "d"), -1)
	})
}

func Test_InArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(InArray(a, "a"), true)
		t.AssertEQ(InArray(a, "b"), true)
		t.AssertEQ(InArray(a, "c"), true)
		t.AssertEQ(InArray(a, "d"), false)
	})
}
