package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/argon2"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("加解密")
	w.Resize(fyne.NewSize(800, 800))

	key := widget.NewMultiLineEntry()
	kc := container.NewVBox(widget.NewLabel("密码："), key)
	content := widget.NewMultiLineEntry()
	cc := container.NewVBox(widget.NewLabel("待处理内容："), content)
	content.SetMinRowsVisible(5)
	encrypt := widget.NewButton("加密", func() {
		handle(key, content, true)
	})
	decrypt := widget.NewButton("解密", func() {
		handle(key, content, false)
	})

	op := container.NewHBox(encrypt, decrypt)
	c := container.NewVBox(kc, cc, op, widget.NewLabel("加解密全在本机进行"))
	w.SetContent(c)
	w.ShowAndRun()
}

func handle(key, content *widget.Entry, encrypt bool) {
	if key.Text == "" {
		key.SetText("请输入密钥")
		return
	}
	salt := sha256.Sum256([]byte(key.Text))
	aeskey := argon2.IDKey([]byte(key.Text), salt[:], 1, 64*1024, 4, 32)

	if content.Text == "" {
		if encrypt {
			content.SetText("请输入待加密内容")
		} else {
			content.SetText("请输入待解密内容")
		}
		return
	}

	c, err := aes.NewCipher(aeskey)
	if err != nil {
		panic(err)
	}
	a, err := cipher.NewGCMWithRandomNonce(c)
	if err != nil {
		panic(err)
	}
	if encrypt {
		content.SetText(base64.StdEncoding.EncodeToString((a.Seal(nil, nil, []byte(content.Text), nil))))
	} else {
		msg, err := base64.StdEncoding.DecodeString(content.Text)
		if err != nil {
			content.SetText(fmt.Sprintf("解密失败：%s", err.Error()))
			return
		}
		b, err := a.Open(nil, nil, msg, nil)
		if err != nil {
			content.SetText(fmt.Sprintf("解密失败：%s", err.Error()))
			return
		}
		content.SetText(string(b))
	}
}

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

//go:embed ..\DingLieSongKeTi\dingliesongtypeface20241217-2.ttf
var ttf []byte

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return fyne.NewStaticResource("ttf", ttf)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
