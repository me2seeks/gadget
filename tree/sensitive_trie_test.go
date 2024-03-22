package tree

import (
	"fmt"
	"testing"
)

func TestT(t *testing.T) {
	sensitiveWords := []string{
		"傻逼",
		"傻叉",
		"垃圾",
		"妈的",
		"sb",
		"牛大大",
	}

	matchContents := []string{
		"你是一个大傻逼，大傻叉",
		"你是傻☺叉",
		"shabi东西",
		"他made东西",
		"什么垃 圾打野，傻逼一样，叫你来开龙不来，SB",
		"正常的内容☺",
	}

	fmt.Println("\n--------- 前缀树匹配敏感词 ---------")
	trieDemo(sensitiveWords, matchContents)

}

// 前缀树匹配敏感词
func trieDemo(sensitiveWords []string, matchContents []string) {

	// 汉字转拼音
	pinyinContents := HansCovertPinyin(sensitiveWords)
	fmt.Println(pinyinContents)

	trie := NewSensitiveTrie()

	trie.AddWords(sensitiveWords)
	trie.AddWords(pinyinContents)

	for _, srcText := range matchContents {
		matchSensitiveWords, replaceText := trie.Match(srcText)
		fmt.Println("srcText        -> ", srcText)
		fmt.Println("replaceText    -> ", replaceText)
		fmt.Println("sensitiveWords -> ", matchSensitiveWords)
		fmt.Println()
	}

	// 动态添加
	trie.AddWord("牛大大")
	content := "今天，牛大大去挑战灰大大了"
	matchSensitiveWords, replaceText := trie.Match(content)
	fmt.Println("srcText        -> ", content)
	fmt.Println("replaceText    -> ", replaceText)
	fmt.Println("sensitiveWords -> ", matchSensitiveWords)
}
