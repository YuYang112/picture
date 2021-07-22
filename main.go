package main

import (
	"context"
	"fmt"
	"github.com/golang/freetype"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"picture/anchor"
	"strconv"
	"strings"
)

//切图时横坐标连续点个数
const (
	THRESHOLD = 3
	ALPHA     = 25
	BTMPATH   = `./image/_btm/`
	XIAWUCHA  = `./xiawucha.ttf`
)

func main() {
	files, dirs, _ := GetFilesAndDirs("./image")

	for _, table := range dirs {
		temp, _, _ := GetFilesAndDirs(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	for _, filepath := range files {
		if err := do(filepath); err != nil {
			fmt.Printf("%s文件处理失败，err：%s\n", filepath, err.Error())
		}

	}
	fmt.Printf("==================图片处理完成=====================\n")

}

func do(filepath string) (err error) {
	err = os.MkdirAll(BTMPATH, 0755)

	if err != nil {
		return
	}
	f, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()
	img, err := png.Decode(f)
	if err != nil {
		return
	}

	maxY := findMaxY(img)
	minY := findMinY(img)
	maxX := findMaxX(img)
	minX := findMinX(img)

	err = cut(minX, maxY, maxX, minY, filepath, img)
	if err != nil {
		return
	}
	return
}

func findMaxY(img image.Image) (maxY int) {
	//从下到上，从左到右
	for y := img.Bounds().Dy(); y >= 0; y-- {
		xNum := 0
		flag := 0
		for x := 3; x < img.Bounds().Dx(); x++ {
			if _, _, _, a := img.At(x, y).RGBA(); uint8(a>>8) > ALPHA {

				if xNum+1 == x {
					flag++
				}
				xNum = x
				if flag > THRESHOLD {
					maxY = y
					return
				}
			} else {
				flag = 0
			}
		}

	}
	return
}
func findMaxX(img image.Image) (maxX int) {
	//从右到左，从上到下
	for x := img.Bounds().Dx(); x >= 0; x-- {
		yNum := 0
		flag := 0
		for y := 0; y < img.Bounds().Dy(); y++ {
			if _, _, _, a := img.At(x, y).RGBA(); uint8(a>>8) > ALPHA {
				if yNum+1 == y {
					flag++
				}
				yNum = y
				if flag > THRESHOLD {
					maxX = x
					return
				}
			} else {
				flag = 0
			}
		}
	}
	return
}

func findMinX(img image.Image) (minX int) {
	//从左到右，从下到上
	for x := 0; x <= img.Bounds().Dx(); x++ {
		yNum := 0
		flag := 0
		for y := img.Bounds().Dy(); y >= 0; y-- {
			if _, _, _, a := img.At(x, y).RGBA(); uint8(a>>8) > ALPHA {
				if yNum-1 == y {
					flag++
				}
				yNum = y
				if flag > THRESHOLD {
					minX = x
					return
				}
			} else {
				flag = 0
			}
		}
	}

	return
}
func findMinY(img image.Image) (minY int) {
	//从上到下，从左到右
	for y := 0; y <= img.Bounds().Dy(); y++ {
		xNum := 0
		flag := 0
		for x := 0; x < img.Bounds().Dx(); x++ {
			if _, _, _, a := img.At(x, y).RGBA(); uint8(a>>8) > ALPHA {
				if xNum+1 == x {
					flag++
				}
				xNum = x
				if flag > THRESHOLD {
					minY = y
					return
				}
			} else {
				flag = 0
			}
		}
	}
	return
}

func cut(minX, maxY, maxX, minY int, filePath string, srcImg image.Image) error {
	rect := image.Rect(0, 0, srcImg.Bounds().Dx(), srcImg.Bounds().Dy())
	rgba := image.NewRGBA(rect)
	//====================================================================================

	fontBytes, err := ioutil.ReadFile(XIAWUCHA)
	if err != nil {
		log.Printf("ioutil.ReadFile error：%s", err.Error())
		return err
	}
	//载入字体数据
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Printf("载入字体失败 error：%s", err.Error())
		return err
	}

	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(100)
	//设置字体
	f.SetFont(font)
	//设置尺寸
	f.SetFontSize(16)
	OffsetY := 20
	pt2OffsetY := 50
	zjOffset := 10
	if srcImg.Bounds().Dx() > 1500 {
		OffsetY = 50
		pt2OffsetY = 70
		zjOffset = 35
		f.SetFontSize(26)
	}
	f.SetClip(rgba.Bounds())
	//设置输出的图片
	f.SetDst(rgba)
	//设置字体颜色(红色)
	f.SetSrc(image.NewUniform(color.RGBA{R: 255, G: 0, B: 0, A: 255}))

	//====================================================================================

	for x := 0; x < srcImg.Bounds().Dx(); x++ {
		for y := 0; y <= srcImg.Bounds().Dy(); y++ {
			r, g, b, a := srcImg.At(x, y).RGBA()
			rgba.SetRGBA(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)})
		}
	}

	//横线
	for x := minX; x <= maxX; x++ {
		rgba.Set(x, maxY+OffsetY, color.RGBA{R: 0, G: 0, B: 0, A: 255})
	}
	//设置横线字体的位置
	pt := freetype.Pt(srcImg.Bounds().Dx()/2, maxY+OffsetY+20)

	_, err = f.DrawString(strconv.Itoa(maxX-minX), pt)
	if err != nil {
		log.Fatal(err)
		return err
	}

	//竖线
	for y := minY; y <= maxY; y++ {
		rgba.Set(minX-10, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})

	}
	//设置竖线字体的位置
	pt2 := freetype.Pt(minX-pt2OffsetY, srcImg.Bounds().Dy()/2)

	_, err = f.DrawString(strconv.Itoa(maxY-minY), pt2)
	if err != nil {
		log.Fatal(err)
		return err
	}
	dir := strings.Split(path.Base(filePath), "\\")
	//===========================================================================================
	//画框
	ctx := context.WithValue(context.Background(), "traceId", "46254")

	bs, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	resp, err := anchor.Send(ctx, anchor.ExteriorUri, bs, dir[2])

	var zx []int
	for _, v := range resp.Result.ExteriorAnchors {
		if v.Id == 3 {
			fmt.Printf("resp===>%+v\n", v)
			zx = append(zx, int(v.Box[0]+(v.Box[2]-v.Box[0])/2))
			//画框
			for x := v.Box[0]; x <= v.Box[2]; x++ {
				rgba.Set(int(x), int(v.Box[1]), color.RGBA{R: 255, G: 255, B: 0, A: 255})
				rgba.Set(int(x), int(v.Box[3]), color.RGBA{R: 255, G: 255, B: 0, A: 255})
			}

			for y := v.Box[1]; y <= v.Box[3]; y++ {
				rgba.Set(int(v.Box[0]), int(y), color.RGBA{R: 255, G: 255, B: 0, A: 255})
				rgba.Set(int(v.Box[2]), int(y), color.RGBA{R: 255, G: 255, B: 0, A: 255})
			}
		}

	}
	if len(zx) >= 2 {

		pt3 := freetype.Pt(srcImg.Bounds().Dx()/2, maxY+OffsetY-zjOffset)

		if zx[1] > zx[0] {
			for x := zx[0]; x <= zx[1]; x++ {
				rgba.Set(x, maxY+OffsetY-zjOffset, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
			_, err = f.DrawString("轴距："+strconv.Itoa(zx[1]-zx[0]), pt3)
		} else {
			for x := zx[1]; x <= zx[0]; x++ {
				rgba.Set(x, maxY+OffsetY-zjOffset, color.RGBA{R:0, G:0, B:0, A:255})
			}
			_, err = f.DrawString("轴距："+strconv.Itoa(zx[0]-zx[1]), pt3)
		}

		if err != nil {
			log.Fatal(err)
			return err
		}

	}

	if err != nil {
		return err
	}
	//===========================================================================================

	_ = os.MkdirAll(BTMPATH+dir[1], 0755)

	distFile, err := os.Create(path.Join(BTMPATH, dir[1], dir[2]))
	if err != nil {
		return err
	}

	defer func() {
		_ = distFile.Close()
	}()
	err = png.Encode(distFile, rgba)
	return err

}

//获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.Name() == "_btm" {
			continue
		}
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			_, _, _ = GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}

	return files, dirs, nil
}
