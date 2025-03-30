package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/argon2"
)

func main() {
	a := app.New()
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
	c := container.NewVBox(kc, cc, op)
	w.SetContent(c)
	w.ShowAndRun()
}

func handle(key, content *widget.Entry, encrypt bool) {
	if key.Text == "" {
		key.SetText("请输入密钥")
		return
	}
	aeskey := argon2.IDKey([]byte(key.Text), nil, 1, 64*1024, 4, 32)

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
