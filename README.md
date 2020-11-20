# minitools
go tools

# myimage
```go
package main

import (
	myImg "day01/minitools/myimage"
	"fmt"
)

func main() {

	img := &myImg.Picture{"1.jpg", nil, nil}
	img.LoadImg()
	fmt.Println(img.GetSize())

	newImg := &myImg.Picture{"3.jpg", nil, nil}
	// img.Copy(newImg)
	// img.Crop(newImg, image.Rect(0, 0, 300, 244))
	// img.ToGray(newImg)
	// img.ColorReverse(newImg)
	// img.HorizontalFlip(newImg)
	// img.VerticalFlip(newImg)
	// img.Rotate(newImg, 45)

	// 普利维特算子(Prewitt operate)
	// img.Filter(newImg, [9]float32{-1, 0, 1, -1, 0, 1, -1, 0, 1})
	// img.Filter(newImg, [9]float32{1, 1, 1, 0, 0, 0, -1, -1, -1})

	// 索贝尔算子（Sobel operator）
	// img.Filter(newImg, [9]float32{-1, 0, 1, -2, 0, 2, -1, 0, 1})
	// img.Filter(newImg, [9]float32{1, 2, 1, 0, 0, 0, -1, -2, -1})

	// 拉普拉斯算子
	// img.Filter(newImg, [9]float32{0, -1, 0, -1, 4, -1, 0, -1, 0})
	// img.Filter(newImg, [9]float32{-1, -1, -1, -1, 8, -1, -1, -1, -1})
	// img.Filter(newImg, [9]float32{0, 1, 0, 1, -4, 1, 0, 1, 0})

	// 均值滤波
	// img.Filter(newImg, [9]float32{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0})

	// 中值滤波
	// img.MedianFilter(newImg, 3)

	// 高斯滤波
	// img.Filter(newImg, [9]float32{1.0 / 16, 2.0 / 16, 1.0 / 16, 2.0 / 16, 4.0 / 16, 2.0 / 16, 1.0 / 16, 2.0 / 16, 1.0 / 16})

	// 改变亮度
	// img.Brightness(newImg, [3]float32{1.2, 1.2, 1.2})

	// img.SaltNoise(newImg, 0.9)
	// img.GaussianNoise(newImg, 1, 0.8)

	// img.GradientImage(newImg, "xy")

	// img.Resize(newImg, 300, 300, "bilinear")

	// fmt.Println(string(img.ImgToBase64()))

	bs64 := myImg.FileToBase64("1.jpg")
	bbb := myImg.Base642buffer(bs64)
	subImg := myImg.BufferToImg(bbb)
	newImg.Img = subImg

	newImg.Save("4.jpg")

}

```

# logging

```go
package main

import (
	logging "day01/minitools/logging"
	"fmt"
	"time"
)

func main() {
	logger, err := logging.NewLogger("debug", "log.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		logger.Debug("debug")
		logger.Debug("debug2")
		logger.Trace("TRACE")
		logger.Info("INFO")
		logger.Warning("WARNING")
		logger.Error("ERROR")
		logger.Fatal("FATAL")

		time.Sleep(time.Second * 1)
	}

}
```