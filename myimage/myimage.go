package myimage

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"
)

/*
图片的常用操作
*/

// Picture struct
type Picture struct {
	ImgPath string
	File    *os.File
	Img     image.Image
}

// LoadImg 加载图片
func (p *Picture) LoadImg() (err error) {
	f, err := os.Open(p.ImgPath)
	if err != nil {
		return
	}
	defer f.Close()
	// p.File = f

	Img, err := jpeg.Decode(f)
	if err != nil {
		return
	}
	p.Img = Img

	return
}

// Close 关闭文件
func (p *Picture) Close() (err error) {
	if p.File != nil {
		p.File.Close()
	}
	return nil
}

// GetSize 获取图片的宽和高
func (p *Picture) GetSize() (int, int) {
	// size := p.Img.Bounds().Max
	size := p.Img.Bounds().Size()
	return size.X, size.Y
}

// Save 图片保存
func (p *Picture) Save(newPath string) (err error) {
	// 保存图像
	f, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer f.Close()
	jpeg.Encode(f, p.Img, &jpeg.Options{100})

	return
}

// Copy 复制图片
func (p *Picture) Copy(p1 *Picture) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			// fmt.Println(Img.At(j, i))
			newImg.SetRGBA(j, i, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	p1.Img = newImg

	return
}

// Clip 像素值裁剪到指定的范围
func Clip(value, min, max float32) uint8 {
	if value < min {
		value = min
	} else if value > max {
		value = max
	}
	return uint8(value)
}

// Crop 按指定大小裁剪
func (p *Picture) Crop(p1 *Picture, r image.Rectangle) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			newImg.SetRGBA(j, i, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}

	// crop
	p1.Img = newImg.SubImage(r)

	return
}

// ToGray 图片灰度化
func (p *Picture) ToGray(p1 *Picture) (err error) {
	w, h := p.GetSize()
	newImg := image.NewGray(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := p.Img.At(j, i).RGBA()
			// newImg.SetGray(j, i, color.Gray{uint8(uint32(0.39*float32(r)+0.5*float32(g)+0.11*float32(b)) >> 8)})
			newImg.SetGray(j, i, color.Gray{Clip(0.39*float32(r>>8)+0.5*float32(g>>8)+0.11*float32(b>>8), float32(0.0), float32(255.0))})
		}
	}
	p1.Img = newImg

	return
}

// ColorReverse 图片像素值反转
func (p *Picture) ColorReverse(p1 *Picture) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			newImg.SetRGBA(j, i, color.RGBA{uint8(255 - r>>8), uint8(255 - g>>8), uint8(255 - b>>8), uint8(255 - a>>8)})
		}
	}
	p1.Img = newImg

	return
}

// HorizontalFlip 水平翻转(左右镜像)
func (p *Picture) HorizontalFlip(p1 *Picture) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			newImg.SetRGBA(w-1-j, i, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	p1.Img = newImg

	return
}

// VerticalFlip 垂直翻转(上下镜像)
func (p *Picture) VerticalFlip(p1 *Picture) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			newImg.SetRGBA(j, h-1-i, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	p1.Img = newImg

	return
}

// U8color 存储uint8类型的颜色值
type U8color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// NewU8color 构造函数
func NewU8color(img image.Image, x, y int) U8color {
	r, g, b, a := img.At(x, y).RGBA()

	// uint32 --> uint8
	rUint8, gUint8, bUint8, aUint8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)

	return U8color{rUint8, gUint8, bUint8, aUint8}
}

// PixelIsZero 判断像素值是否都为0
func (u U8color) PixelIsZero() bool {
	if u.Red == uint8(0) && u.Green == uint8(0) && u.Blue == uint8(0) {
		return true
	}
	return false
}

// BilinearInterpolation 双线性插值
func BilinearInterpolation(img image.Image, w, h int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			c11 := NewU8color(img, j, i)
			if c11.PixelIsZero() {
				// 只使用临近4个点做插值
				c01 := NewU8color(img, j-1, i)
				c21 := NewU8color(img, j+1, i)
				c10 := NewU8color(img, j, i-1)
				c12 := NewU8color(img, j, i+1)

				newR := Clip((float32(c01.Red)+float32(c21.Red)+float32(c10.Red)+float32(c12.Red))/float32(4.0), float32(0.0), float32(255.0))
				newG := Clip((float32(c01.Green)+float32(c21.Green)+float32(c10.Green)+float32(c12.Green))/float32(4.0), float32(0.0), float32(255.0))
				newB := Clip((float32(c01.Blue)+float32(c21.Blue)+float32(c10.Blue)+float32(c12.Blue))/float32(4.0), float32(0.0), float32(255.0))
				newA := Clip((float32(c01.Alpha)+float32(c21.Alpha)+float32(c10.Alpha)+float32(c12.Alpha))/float32(4.0), float32(0.0), float32(255.0))

				newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
			} else {
				newImg.SetRGBA(j, i, color.RGBA{c11.Red, c11.Green, c11.Blue, c11.Alpha})
			}
		}

	}

	return newImg
}

// Rotate 旋转
func (p *Picture) Rotate(p1 *Picture, angle float64) (err error) {
	angle = angle / 180.0 * math.Pi
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	cx, cy := float64(w)/2.0, float64(h)/2.0
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.Img.At(j, i).RGBA()
			// newImg.SetRGBA(j, h-1-i, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			// 中心点旋转
			j1 := float64(j) - cx
			i1 := float64(i) - cy
			j1, i1 = j1*math.Cos(angle)-i1*math.Sin(angle), j1*math.Sin(angle)+i1*math.Cos(angle)
			// 回到原中心点
			j1 += cx
			i1 += cy
			// 判断是否越界
			j2, i2 := int(j1), int(i1)
			if j2 < 0 || j2 >= w || i2 < 0 || i2 >= h {
				continue
			}
			newImg.SetRGBA(j2, i2, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	// p1.Img = newImg
	// 使用插值
	p1.Img = BilinearInterpolation(newImg, w, h)

	return
}

// Filter 滤波
func (p *Picture) Filter(p1 *Picture, arr [9]float32) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	img := p.Img
	// 修改像素值
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			c11 := NewU8color(img, j, i)
			// 距离为1
			c01 := NewU8color(img, j-1, i)
			c21 := NewU8color(img, j+1, i)
			c10 := NewU8color(img, j, i-1)
			c12 := NewU8color(img, j, i+1)
			// 距离为2
			c00 := NewU8color(img, j-1, i-1)
			c02 := NewU8color(img, j-1, i+1)
			c20 := NewU8color(img, j+1, i-1)
			c22 := NewU8color(img, j+1, i+1)

			newR := Clip(float32(c00.Red)*arr[0]+float32(c01.Red)*arr[1]+float32(c02.Red)*arr[2]+
				float32(c10.Red)*arr[3]+float32(c11.Red)*arr[4]+float32(c12.Red)*arr[5]+
				float32(c20.Red)*arr[6]+float32(c21.Red)*arr[7]+float32(c22.Red)*arr[8], float32(0.0), float32(255.0))

			newG := Clip(float32(c00.Green)*arr[0]+float32(c01.Green)*arr[1]+float32(c02.Green)*arr[2]+
				float32(c10.Green)*arr[3]+float32(c11.Green)*arr[4]+float32(c12.Green)*arr[5]+
				float32(c20.Green)*arr[6]+float32(c21.Green)*arr[7]+float32(c22.Green)*arr[8], float32(0.0), float32(255.0))

			newB := Clip(float32(c00.Blue)*arr[0]+float32(c01.Blue)*arr[1]+float32(c02.Blue)*arr[2]+
				float32(c10.Blue)*arr[3]+float32(c11.Blue)*arr[4]+float32(c12.Blue)*arr[5]+
				float32(c20.Blue)*arr[6]+float32(c21.Blue)*arr[7]+float32(c22.Blue)*arr[8], float32(0.0), float32(255.0))

			newA := Clip(float32(c00.Alpha)*arr[0]+float32(c01.Alpha)*arr[1]+float32(c02.Alpha)*arr[2]+
				float32(c10.Alpha)*arr[3]+float32(c11.Alpha)*arr[4]+float32(c12.Alpha)*arr[5]+
				float32(c20.Alpha)*arr[6]+float32(c21.Alpha)*arr[7]+float32(c22.Alpha)*arr[8], float32(0.0), float32(255.0))

			newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})

		}

	}

	p1.Img = newImg
	return
}

// SortedU8colorSlice ...
func SortedU8colorSlice(tmp []U8color, size int, center int) color.RGBA {
	newR := make([]uint8, 0, size)
	newG := make([]uint8, 0, size)
	newB := make([]uint8, 0, size)
	newA := make([]uint8, 0, size)
	for i := 0; i < size; i++ {
		// newR[i] = tmp[i].Red
		// newG[i] = tmp[i].Green
		// newB[i] = tmp[i].Blue
		// newA[i] = tmp[i].Alpha

		newR = append(newR, tmp[i].Red)
		newG = append(newG, tmp[i].Green)
		newB = append(newB, tmp[i].Blue)
		newA = append(newA, tmp[i].Alpha)
	}

	// 排序
	sort.Slice(newR, func(i, j int) bool {
		return newR[i] < newR[j]
	})
	sort.Slice(newG, func(i, j int) bool {
		return newG[i] < newG[j]
	})
	sort.Slice(newB, func(i, j int) bool {
		return newB[i] < newB[j]
	})
	sort.Slice(newA, func(i, j int) bool {
		return newA[i] < newA[j]
	})

	return color.RGBA{newR[center], newG[center], newB[center], newA[center]}
}

// MedianFilter 中值滤波 默认 3x3 为例
func (p *Picture) MedianFilter(p1 *Picture, ksize int) (err error) {
	img := p.Img
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	pad := ksize / 2

	// 修改像素值
	for i := pad; i < h-pad; i++ {
		for j := pad; j < w-pad; j++ {
			tmp := make([]U8color, 0, ksize*ksize)
			for ky := 0; ky < ksize; ky++ {
				for kx := 0; kx < ksize; kx++ {
					// tmp[kx+ky*ksize] = NewU8color(img, j-pad+kx, i-pad+ky)
					tmp = append(tmp, NewU8color(img, j-pad+kx, i-pad+ky))
				}
			}
			newImg.SetRGBA(j, i, SortedU8colorSlice(tmp, ksize*ksize, pad))
		}

	}

	p1.Img = newImg
	return
}

// Brightness 改变亮度
func (p *Picture) Brightness(p1 *Picture, arr [3]float32) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	img := p.Img
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			c11 := NewU8color(img, j, i)
			newR := Clip(float32(c11.Red)*arr[0], 0, 255)
			newG := Clip(float32(c11.Green)*arr[1], 0, 255)
			newB := Clip(float32(c11.Blue)*arr[2], 0, 255)
			newA := c11.Alpha

			newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})

		}
	}
	p1.Img = newImg
	return
}

// SaltNoise 椒盐噪声
func (p *Picture) SaltNoise(p1 *Picture, snr float32) (err error) {
	// snr 信噪比
	w, h := p.GetSize()
	noiseSize := int(float32(w*h) * (1 - snr))
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	img := p.Img
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			c11 := NewU8color(img, j, i)
			newImg.SetRGBA(j, i, color.RGBA{c11.Red, c11.Green, c11.Blue, c11.Alpha})
		}
	}

	// 设置噪声
	for k := 0; k < noiseSize; k++ {
		// 随机获取 某个点
		x := rand.Intn(w)
		y := rand.Intn(h)
		// 增加噪声
		newImg.SetRGBA(x, y, color.RGBA{uint8(0), uint8(0), uint8(0), uint8(0)})
	}

	p1.Img = newImg
	return
}

// GaussianNoise 高斯噪声
func (p *Picture) GaussianNoise(p1 *Picture, mu, sigma float64) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	img := p.Img
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			c11 := NewU8color(img, j, i)
			u1 := rand.Float64()
			u2 := rand.Float64()
			z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2*math.Pi*u2)
			// z1 := math.Sqrt(-2.0*math.Log(u1)) * math.Sin(2*math.Pi*u2)
			noise := (z0*sigma + mu) * 32

			newR := Clip(float32(noise)+float32(c11.Red), 0, 255)
			newG := Clip(float32(noise)+float32(c11.Green), 0, 255)
			newB := Clip(float32(noise)+float32(c11.Blue), 0, 255)
			newA := Clip(float32(noise)+float32(c11.Alpha), 0, 255)

			newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
		}
	}

	p1.Img = newImg
	return
}

// GradientImage 梯度图像
func (p *Picture) GradientImage(p1 *Picture, mode string) (err error) {
	w, h := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))

	img := p.Img
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			c11 := NewU8color(img, j, i)
			if "x" == strings.ToLower(mode) { // 水平方向梯度图
				if j < w-1 {
					c21 := NewU8color(img, j+1, i)
					newR := Clip(float32(c21.Red)-float32(c11.Red), 0, 255)
					newG := Clip(float32(c21.Green)-float32(c11.Green), 0, 255)
					newB := Clip(float32(c21.Blue)-float32(c11.Blue), 0, 255)
					newA := Clip(float32(c21.Alpha)-float32(c11.Alpha), 0, 255)
					newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
				}
				// } else {
				// 	newImg.SetRGBA(j, i, color.RGBA{c11.Red, c11.Green, c11.Blue, c11.Alpha})
				// }
			} else if "y" == strings.ToLower(mode) { // 垂直方向梯度图
				if i < h-1 {
					c12 := NewU8color(img, j, i+1)
					newR := Clip(float32(c12.Red)-float32(c11.Red), 0, 255)
					newG := Clip(float32(c12.Green)-float32(c11.Green), 0, 255)
					newB := Clip(float32(c12.Blue)-float32(c11.Blue), 0, 255)
					newA := Clip(float32(c12.Alpha)-float32(c11.Alpha), 0, 255)
					newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
				}
				// } else {
				// 	newImg.SetRGBA(j, i, color.RGBA{c11.Red, c11.Green, c11.Blue, c11.Alpha})
				// }
			} else { // 垂直方向+ 水平方向
				if i < h-1 && j < w-1 {
					c12 := NewU8color(img, j, i+1)
					c21 := NewU8color(img, j+1, i)

					newR := Clip(float32(c21.Red)+float32(c12.Red)-float32(c11.Red)*2, 0, 255)
					newG := Clip(float32(c21.Green)+float32(c12.Green)-float32(c11.Green)*2, 0, 255)
					newB := Clip(float32(c21.Blue)+float32(c12.Blue)-float32(c11.Blue)*2, 0, 255)
					newA := Clip(float32(c21.Alpha)+float32(c12.Alpha)-float32(c11.Alpha)*2, 0, 255)

					newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
				}
			}

		}
	}

	p1.Img = newImg
	return
}

// Resize ... 使用双线性插值
func (p *Picture) Resize(p1 *Picture, w, h int, mode string) (err error) {
	imgW, imgH := p.GetSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	img := p.Img
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			// 反算出对应与原图的坐标
			x := (float64(j)+0.5)*float64(imgW)/float64(w) - 0.5
			y := (float64(i)+0.5)*float64(imgH)/float64(h) - 0.5

			x1 := int(math.Floor(x)) // 向下取整
			y1 := int(math.Floor(y)) // 向下取整

			x2 := int(math.Ceil(x)) // # 向上取整
			y2 := int(math.Ceil(y)) // # 向上取整

			if x1 < 0 {
				x1 = 0
			} else if x2 >= imgW {
				x2 = imgW - 1
			}

			if y1 < 0 {
				y1 = 0
			} else if y2 >= imgH {
				y2 = imgH - 1
			}

			if mode == "bilinear" { // 双线性插值

				c11 := NewU8color(img, x1, y1)
				c12 := NewU8color(img, x1, y2)
				c21 := NewU8color(img, x2, y1)
				c22 := NewU8color(img, x2, y2)

				newR := Clip(float32(float64(c11.Red)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c12.Red)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y2)-y))+
					float64(c21.Red)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c22.Red)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y2)-y))), 0, 255)

				newG := Clip(float32(float64(c11.Green)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c12.Green)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y2)-y))+
					float64(c21.Green)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c22.Green)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y2)-y))), 0, 255)

				newB := Clip(float32(float64(c11.Blue)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c12.Blue)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y2)-y))+
					float64(c21.Blue)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c22.Blue)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y2)-y))), 0, 255)

				newA := Clip(float32(float64(c11.Alpha)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c12.Alpha)*(1-math.Abs(float64(x1)-x))*(1-math.Abs(float64(y2)-y))+
					float64(c21.Alpha)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y1)-y))+
					float64(c22.Alpha)*(1-math.Abs(float64(x2)-x))*(1-math.Abs(float64(y2)-y))), 0, 255)

				newImg.SetRGBA(j, i, color.RGBA{newR, newG, newB, newA})
			} else if mode == "nearest" { // 最近邻插值
				var newX, newY = 0, 0
				if math.Abs(float64(x1)-x) <= math.Abs(float64(x2)-x) {
					newX = x1
				} else {
					newX = x2
				}

				if math.Abs(float64(y1)-y) <= math.Abs(float64(y2)-y) {
					newY = y1
				} else {
					newY = y2
				}

				c11 := NewU8color(img, newX, newY)
				newImg.SetRGBA(j, i, color.RGBA{c11.Red, c11.Green, c11.Blue, c11.Alpha})
			} else {
				err = errors.New("mode only nearest or bilinear")
				return
			}
		}
	}
	p1.Img = newImg
	return
}

// ImgToBase64 image.Image转成 base64
func (p *Picture) ImgToBase64() string {
	// 开辟一个新的空buff
	emptyBuff := bytes.NewBuffer(nil)
	// img写入到buff
	jpeg.Encode(emptyBuff, p.Img, nil)
	// 开辟存储空间
	dist := make([]byte, 50000)
	// buff转成base64
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	// 输出图片base64(type = []byte)
	// fmt.Println(string(dist))
	// buffer输出到jpg文件中（不做处理，直接写到文件）
	// _ = ioutil.WriteFile("./base64pic.txt", dist, 0666)

	return string(dist)
}

// FileToBase64  图片路径直接转成 base64
func FileToBase64(imgPath string) string {
	ff, _ := ioutil.ReadFile(imgPath)       // 读取整个文件
	bufstore := make([]byte, 5000000)       //数据缓存
	base64.StdEncoding.Encode(bufstore, ff) // 文件转base64
	// _ = ioutil.WriteFile("./output2.jpg.txt", bufstore, 0666) //直接写入到文件就ok。

	return string(bufstore)
}

// Base642File  base64直接转成保存成图片
func Base642File(datasource, imgPath string) (err error) {
	// 成图片文件并把文件写入到buffer
	ddd, _ := base64.StdEncoding.DecodeString(datasource)
	// buffer输出到jpg文件中（不做处理，直接写到文件）
	err = ioutil.WriteFile(imgPath, ddd, 0666)
	return
}

// Base642buffer  base64直接转成保存成图片
func Base642buffer(datasource string) *bytes.Buffer {
	// 成图片文件并把文件写入到buffer
	ddd, _ := base64.StdEncoding.DecodeString(datasource)
	// 必须加一个buffer 不然没有read方法就会报错
	bbb := bytes.NewBuffer(ddd)
	return bbb
}

// BufferToImg  base64直接转成保存成图片
func BufferToImg(bbb *bytes.Buffer) image.Image {
	m, _, _ := image.Decode(bbb) // 图片文件解码
	// rgbImg := m.(*image.YCbCr)
	// subImg := rgbImg.SubImage(image.Rect(0, 0, w, h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1

	return m
}
