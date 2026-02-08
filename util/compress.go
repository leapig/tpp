package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/nfnt/resize"
)

func CompressBase64(base64Str string, targetSize int64) (string, error) {
	maxTargetSize := targetSize * 1000
	size := getSize(base64Str)
	if size <= maxTargetSize {
		// 小于最小压缩尺寸，不压缩
		return base64Str, nil
	} else {
		// 大于最大压缩尺寸，压缩到最大尺寸
		return compress(resizeImage(base64Str), maxTargetSize), nil
	}
}

func compress(str string, maxTargetSize int64) string {
	var (
		quality     = 100
		currentStr  = str
		minQuality  = 30
		maxAttempts = 10 // 防止无限循环的安全阀
	)

	for attempt := 0; attempt < maxAttempts; attempt++ {
		size := getSize(currentStr)
		if size <= maxTargetSize || quality <= minQuality {
			break
		}

		// 动态调整质量参数
		ratio := float64(maxTargetSize) / float64(size)
		quality = int(float64(quality) * ratio)
		if quality < minQuality {
			quality = minQuality
		}

		// 图像处理流程
		body, err := base64.StdEncoding.DecodeString(currentStr)
		if err != nil {
			fmt.Println("base64 decode error:", err)
			return currentStr
		}

		img, _, err := image.Decode(bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("image decode error:", err)
			return currentStr
		}

		var buffer bytes.Buffer
		if err := jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality}); err != nil {
			fmt.Println("jpeg encode error:", err)
			return currentStr
		}

		currentStr = base64.StdEncoding.EncodeToString(buffer.Bytes())
		quality = int(float64(quality) * 0.9) // 渐进式质量下降
	}
	return currentStr

	//
	//size := getSize(str)
	//fmt.Println("compress file size is:", size)
	//if size <= maxTargetSize {
	//	// 小于最小压缩尺寸，不压缩
	//	return str
	//} else {
	//	// 大于最大压缩尺寸，压缩到最大尺寸
	//	radio := big.NewFloat(0)
	//	radio.Quo(big.NewFloat(float64(maxTargetSize)), big.NewFloat(float64(size)))
	//	result := big.NewFloat(0)
	//	//result.Add(radio, big.NewFloat(0.02))
	//	ra := result.Text('f', 2)
	//	floatValue, _ := strconv.ParseFloat(ra, 64)
	//	quality := int(floatValue * 100)
	//	body, _ := base64.StdEncoding.DecodeString(str)
	//	img, _, e := image.Decode(bytes.NewBuffer(body))
	//	if e != nil {
	//		fmt.Println("compress error is:", e)
	//	}
	//	var buffer bytes.Buffer
	//	if quality > 100 {
	//		quality = 100
	//	} else if quality < 30 {
	//		quality = 30
	//	}
	//	_ = jpeg.Encode(&buffer, img, &jpeg.Options{Quality: quality})
	//	res := base64.StdEncoding.EncodeToString(buffer.Bytes())
	//	return compress(res, maxTargetSize)
	//}
}

func resizeImage(str string) string {
	body, _ := base64.StdEncoding.DecodeString(str)
	img, _, e := image.Decode(bytes.NewBuffer(body))
	if e != nil {
		fmt.Println("resizeImage error:", e)
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	if width <= 420 && height <= 420 {
		return str
	}
	var targetWidth, targetHeight uint
	if width >= height {
		radio := height / 420
		targetWidth = uint(width / radio)
		targetHeight = 420
	} else {
		radio := width / 420
		targetWidth = 420
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
