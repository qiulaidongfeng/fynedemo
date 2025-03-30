package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	_ "embed"
)

//go:embed base.cn.json
var cn []byte

//go:embed word.txt
var word string

var wordslist []Word

type Word struct {
	v string
}

// parserWord 解析以换行和空格分隔的单词列表
func parserWord(s string) {
	s = strings.Replace(s, "\n", " ", -1)
	v := strings.Split(s, " ")
	for _, w := range v {
		if strings.Contains(w, " ") || len(w) == 0 {
			panic(w)
		}
		wordslist = append(wordslist, Word{v: w})
	}
}

func init() {
	parserWord(word)
}

func main() {
	a := app.New()
	lang.AddTranslationsForLocale(cn, fyne.Locale("zh-cn"))

	w := a.NewWindow("随机选单词")
	w.Resize(fyne.NewSize(800, 800))

	c := container.New(layout.NewCustomPaddedVBoxLayout(20), makeUI()...)
	w.SetContent(c)
	w.ShowAndRun()
}

func makeUI() (ret []fyne.CanvasObject) {
	wordlist := widget.NewTextGrid()
	wordlist.SetText(word)

	num := widget.NewEntry()
	num.Resize(fyne.NewSize(20, 100))
	numAll := container.NewVBox(widget.NewLabel("随机选择单词数量："), num)

	result := widget.NewTextGrid()
	r := container.NewScroll(result)
	r.SetMinSize(fyne.NewSize(200, 200))
	resultAll := container.New(layout.NewCustomPaddedVBoxLayout(20), widget.NewLabel("随机选择单词结果："), r)

	Rand := widget.NewButton("随机选单词", func() {
		i, err := strconv.Atoi(num.Text)
		if err != nil {
			num.SetText("输入的数字有误，请重新输入")
			return
		}
		var buf strings.Builder
		for r := range i {
			buf.WriteString(fmt.Sprintf("%d\t", r+1))
			index := rand.Uint32N(uint32(len(wordslist) - 1))
			buf.WriteString(wordslist[index].v)
			if r < i-1 {
				buf.WriteString("\n")
			}
		}
		result.SetText(buf.String())
	})
	x := container.NewScroll(wordlist)
	x.SetMinSize(fyne.NewSize(200, 200))
	ret = append(ret, x, numAll, Rand, resultAll)
	return
}
