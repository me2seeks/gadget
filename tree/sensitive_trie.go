package tree

import (
	"regexp"
	"strings"

	"github.com/Chain-Zhang/pinyin"
)

type SensitiveTrie struct {
	replaceChar rune
	root        *TrieNode
}

type TrieNode struct {
	childMap map[rune]*TrieNode
	isEnd    bool
	data     string
}

func NewSensitiveTrie() *SensitiveTrie {
	return &SensitiveTrie{
		replaceChar: '*',
		root:        &TrieNode{isEnd: false},
	}
}

func (t *TrieNode) AddChild(c rune) *TrieNode {
	if t.childMap == nil {
		t.childMap = make(map[rune]*TrieNode)
	}
	if trieNode, ok := t.childMap[c]; ok {
		return trieNode
	} else {
		t.childMap[c] = &TrieNode{
			childMap: nil,
			isEnd:    false,
		}
		return t.childMap[c]
	}
}

func (t *TrieNode) FindChild(c rune) *TrieNode {
	if t.childMap == nil {
		return nil
	}
	if trieNode, ok := t.childMap[c]; ok {
		return trieNode
	}
	return nil
}

func (st *SensitiveTrie) replaceRune(chars []rune, begin int, end int) {
	for i := begin; i < end; i++ {
		chars[i] = st.replaceChar
	}
}

func (st SensitiveTrie) AddWord(word string) {
	trieNode := st.root
	for _, charInt := range word {
		// 添加敏感词到前缀树中
		trieNode = trieNode.AddChild(charInt)
	}
	trieNode.isEnd = true
	trieNode.data = word
}

// AddWords 批量添加敏感词
func (st *SensitiveTrie) AddWords(sensitiveWords []string) {
	for _, sensitiveWord := range sensitiveWords {
		sensitiveWord := strings.ReplaceAll(sensitiveWord, " ", "")
		st.AddWord(sensitiveWord)
	}
}

func (st *SensitiveTrie) Match(text string) (sensitiveWords []string, replaceText string) {
	if st.root == nil {
		return nil, text
	}
	filteredText := st.FilterSpecialChar(text)
	sensitiveMap := make(map[string]*struct{})

	textChars := []rune(filteredText)
	textCharsCopy := make([]rune, len(textChars))
	copy(textCharsCopy, textChars)
	for i, textLen := 0, len(textChars); i < textLen; i++ {
		trieNode := st.root.FindChild(textChars[i])
		if trieNode == nil {
			continue
		}
		// 匹配到了敏感词的前缀，从后一个位置继续
		j := i + 1
		for ; j < textLen && trieNode != nil; j++ {
			if trieNode.isEnd {
				if _, ok := sensitiveMap[trieNode.data]; !ok {
					sensitiveWords = append(sensitiveWords, trieNode.data)
				}
				st.replaceRune(textCharsCopy, i, j)
			}
			trieNode = trieNode.FindChild(textChars[j])
		}
		// 文本尾部命中敏感词情况
		if j == textLen && trieNode != nil && trieNode.isEnd {
			if _, ok := sensitiveMap[trieNode.data]; !ok {
				sensitiveWords = append(sensitiveWords, trieNode.data)
			}
			sensitiveMap[trieNode.data] = nil
			st.replaceRune(textCharsCopy, i, textLen)
		}
	}

	if len(sensitiveWords) > 0 {
		// 有敏感词
		replaceText = string(textCharsCopy)
	} else {
		// 没有则返回原来的文本
		replaceText = text
	}

	return sensitiveWords, replaceText
}

func (st *SensitiveTrie) FilterSpecialChar(text string) string {
	text = strings.ToLower(text)
	text = strings.Replace(text, " ", "", -1) // 去除空格

	// 过滤除中英文及数字以外的其他字符
	otherCharReg := regexp.MustCompile("[^\u4e00-\u9fa5a-zA-Z0-9]")
	text = otherCharReg.ReplaceAllString(text, "")
	return text
}

func HansCovertPinyin(contents []string) []string {
	pinyinContents := make([]string, 0)
	for _, content := range contents {
		chineseReg := regexp.MustCompile("[\u4e00-\u9fa5]")
		if !chineseReg.Match([]byte(content)) {
			continue
		}
		// 只有中文才转
		pin := pinyin.New(content)
		pinStr, err := pin.Convert()
		println(content, "->", pinStr)
		if err == nil {
			pinyinContents = append(pinyinContents, pinStr)
		}
	}
	return pinyinContents
}
