package util

import (
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

// QrCode 生成指定内容的二维码图像，并返回其PNG格式的字节数据。
//
// 参数:
//   - content: 需要编码为二维码的字符串内容。
//
// 返回值:
//   - []byte: 生成的二维码图像的PNG格式字节数据。
//   - error: 如果生成过程中出现错误，则返回错误信息；否则返回nil。
func QrCode(content string) ([]byte, error) {
	// 将内容编码为二维码的位图
	cs, _ := qr.Encode(content, qr.H, qr.Auto)

	// 将二维码位图缩放为512x512像素
	qrCode, err := barcode.Scale(cs, 512, 512)
	if err != nil {
		return nil, err
	}

	// 创建一个缓冲区用于存储PNG图像数据
	buff := new(bytes.Buffer)

	// 将二维码图像编码为PNG格式并写入缓冲区
	if ee := png.Encode(buff, qrCode); ee != nil {
		return nil, ee
	} else {
		// 返回缓冲区中的PNG字节数据
		return buff.Bytes(), nil
	}
}
