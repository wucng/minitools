package minitools

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

/*
图片的常用操作
*/

// Img struct
type Img struct {
	imgPath string
	file    *os.File
	img     image.Image
}

// LoadImg 加载图片
func (p *Img) LoadImg() (err error) {
	f, err := os.Open(p.imgPath)
	if err != nil {
		return
	}
	defer f.Close()
	// p.file = f

	img, err := jpeg.Decode(f)
	if err != nil {
		return
	}
	p.img = img

	return
}

// Close 关闭文件
func (p *Img) Close() (err error) {
	if p.file != nil {
		p.file.Close()
	}
	return nil
}

// GetSize 获取图片的宽和高
func (p *Img) GetSize() (int, int) {
	// size := p.img.Bounds().Max
	size := p.img.Bounds().Size()
	return size.X, size.Y
}

// Save 图片保存
func (p *Img) Save(newPath string) (err error) {
	// 保存图像
	f, err := os.Create(newPath)
	if err != nil {
		return
	}
	defer f.Close()
	jpeg.Encode(f, p.img, &jpeg.Options{100})

	return
}

// Copy 复制图片
func (p *Img) Copy(p1 *Img) (err error) {
	w, h := p.getSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.img.At(j, i).RGBA()
			// fmt.Println(img.At(j, i))
			newImg.SetRGBA(j, i, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	p1.img = newImg

	return
}

// Crop 按指定大小裁剪
func (p *Img) Crop(p1 *Img, r image.Rectangle) (err error) {
	w, h := p.getSize()
	newImg := image.NewRGBA(image.Rect(0, 0, w, h))
	// 修改像素值
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			r, g, b, a := p.img.At(j, i).RGBA()
			// fmt.Println(img.At(j, i))
			newImg.SetRGBA(j, i, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	// crop
	p1.img = newImg.SubImage(r)

	return
}
