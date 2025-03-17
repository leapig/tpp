package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

func CompressBase64(base64Str string, minTargetSize, maxTargetSize int64) (string, error) {
	maxTargetSize = maxTargetSize * 1000
	minTargetSize = minTargetSize * 1000
	size := getSize(base64Str)
	if size <= maxTargetSize {
		// 小于最小压缩尺寸，不压缩
		return base64Str, nil
	} else {
		// 大于最大压缩尺寸，压缩到最大尺寸
		return compress(resizeImage(base64Str), minTargetSize, maxTargetSize), nil
	}
}

func compress(str string, minTargetSize, maxTargetSize int64) string {
	size := getSize(str)
	fmt.Println("compress file size is:", size)
	if size <= maxTargetSize {
		// 小于最小压缩尺寸，不压缩
		return str
	} else {
		// 大于最大压缩尺寸，压缩到最大尺寸
		radio := big.NewFloat(0)
		radio.Quo(big.NewFloat(float64((maxTargetSize+minTargetSize)/2)), big.NewFloat(float64(size)))
		result := big.NewFloat(0)
		//result.Add(radio, big.NewFloat(0.02))
		ra := result.Text('f', 2)
		floatValue, _ := strconv.ParseFloat(ra, 64)
		quality := int(floatValue * 100)
		body, _ := base64.StdEncoding.DecodeString(str)
		img, _, e := image.Decode(bytes.NewBuffer(body))
		if e != nil {
			fmt.Println("compress error is:", e)
		}
		var buffer bytes.Buffer
		if quality > 100 {
			quality = 100
		} else if quality < 30 {
			quality = 30
		}
		_ = jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality})
		res := base64.StdEncoding.EncodeToString(buffer.Bytes())
		return compress(res, minTargetSize, maxTargetSize)
	}
}

func resizeImage(str string) string {
	body, _ := base64.StdEncoding.DecodeString(str)
	img, _, e := image.Decode(bytes.NewBuffer(body))
	if e != nil {
		fmt.Println("resizeImage error:", e)
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	var targetWidth, targetHeight uint
	if width >= height {
		radio := height / 480
		targetWidth = uint(width / radio)
		targetHeight = 480
	} else {
		radio := width / 480
		targetWidth = 480
		targetHeight = uint(height / radio)
	}
	targetImg := resize.Resize(targetWidth, targetHeight, img, resize.Lanczos3)
	buffer := bytes.NewBuffer(nil)
	_ = png.Encode(buffer, targetImg)
	encoded := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return encoded
}

func getSize(str string) int64 {
	decodedData, _ := base64.StdEncoding.DecodeString(str)
	return int64(len(decodedData))
}

// RandomStrByNum 获取指定位数的数字随机数
func randomStrByNum(n int) string {
	rand.Seed(time.Now().UnixNano())
	number := make([]byte, n)
	number[0] = byte(rand.Intn(9)) + '1' // 确保第一位不为0
	for i := 1; i < n; i++ {
		number[i] = byte(rand.Intn(10)) + '0'
	}
	return string(number)
}
