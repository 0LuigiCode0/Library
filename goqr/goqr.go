package goqr

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

var code = "https://x-cluster.com&key=ijgriufhwunecb23fdfgsfgrgfrgertgerherg"

var ascii = []int{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	36, 0, 0, 0, 37, 38, 0, 0, 0, 0, 39, 40, 0, 41, 42, 43,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 44, 0, 0, 0, 0, 0,
	0, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 0, 0, 0, 0, 0,
	0, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

var maxDataM = []int{
	128, 224, 352, 512, 688, 864, 992, 1232, 1456, 1728,
	2032, 2320, 2672, 2920, 3320, 3624, 4056, 4504, 5016, 5352,
	5712, 6256, 6880, 7312, 8000, 8496, 9024, 9544, 10136, 10984,
	11640, 12328, 13048, 13800, 14496, 15312, 15936, 16816, 17728, 18672,
}
var maxDataH = []int{
	72, 128, 208, 288, 368, 480, 528, 688, 800, 976,
	1120, 1264, 1440, 1576, 1784, 2024, 2264, 2504, 2728, 3080,
	3248, 3536, 3712, 4112, 4304, 4768, 5024, 5288, 5608, 5960,
	6344, 6760, 7208, 7688, 7888, 8432, 8768, 9136, 9776, 10208,
}

var blocksM = []int{
	1, 1, 1, 2, 2, 4, 4, 4, 5, 5,
	5, 8, 9, 9, 10, 10, 11, 13, 14, 16,
	17, 17, 18, 20, 21, 23, 25, 26, 28, 29,
	31, 33, 35, 37, 38, 40, 43, 45, 47, 49,
}
var blocksH = []int{
	1, 1, 2, 4, 4, 4, 5, 6, 8, 8,
	11, 11, 16, 16, 18, 16, 19, 21, 25, 25,
	25, 34, 30, 32, 35, 37, 40, 42, 45, 48,
	51, 54, 57, 60, 63, 66, 70, 74, 77, 81,
}

var byteCorectM = []int{
	10, 16, 26, 18, 24, 16, 18, 22, 22, 26,
	30, 22, 22, 24, 24, 28, 28, 26, 26, 26,
	26, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
}
var byteCorectH = []int{
	17, 28, 22, 16, 22, 28, 26, 26, 24, 28,
	24, 28, 22, 24, 24, 30, 28, 28, 26, 28,
	30, 24, 30, 30, 30, 30, 30, 30, 30, 30,
	30, 30, 30, 30, 30, 30, 30, 30, 30, 30,
}

var polinom = map[int][]int{
	7:  {87, 229, 146, 149, 238, 102, 21},
	10: {251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
	13: {74, 152, 176, 100, 86, 100, 106, 104, 130, 218, 206, 140, 78},
	15: {8, 183, 61, 91, 202, 37, 51, 58, 58, 237, 140, 124, 5, 99, 105},
	16: {120, 104, 107, 109, 102, 161, 76, 3, 91, 191, 147, 169, 182, 194, 225, 120},
	17: {43, 139, 206, 78, 43, 239, 123, 206, 214, 147, 24, 99, 150, 39, 243, 163, 136},
	18: {215, 234, 158, 94, 184, 97, 118, 170, 79, 187, 152, 148, 252, 179, 5, 98, 96, 153},
	20: {17, 60, 79, 50, 61, 163, 26, 187, 202, 180, 221, 225, 83, 239, 156, 164, 212, 212, 188, 190},
	22: {210, 171, 247, 242, 93, 230, 14, 109, 221, 53, 200, 74, 8, 172, 98, 80, 219, 134, 160, 105, 165, 231},
	24: {229, 121, 135, 48, 211, 117, 251, 126, 159, 180, 169, 152, 192, 226, 228, 218, 111, 0, 117, 232, 87, 96, 227, 21},
	26: {173, 125, 158, 2, 103, 182, 118, 17, 145, 201, 111, 28, 165, 53, 161, 21, 245, 142, 13, 102, 48, 227, 153, 145, 218, 70},
	28: {168, 223, 200, 104, 224, 234, 108, 180, 110, 190, 195, 147, 205, 27, 232, 201, 21, 43, 245, 87, 42, 195, 212, 119, 242, 37, 9, 123},
	30: {41, 173, 145, 152, 216, 31, 179, 182, 50, 48, 110, 86, 239, 96, 222, 125, 42, 173, 226, 193, 224, 130, 156, 37, 251, 216, 238, 40, 192, 180},
}

var fieldGalua = []int{
	1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38,
	76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192,
	157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35,
	70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161,
	95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240,
	253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226,
	217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206,
	129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204,
	133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84,
	168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115,
	230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255,
	227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65,
	130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166,
	81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9,
	18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22,
	44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
}

var reversefieldGalua = []int{
	0, 0, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75,
	4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113,
	5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69,
	29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166,
	6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136,
	54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64,
	30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61,
	202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87,
	7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24,
	227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46,
	55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97,
	242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162,
	31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246,
	108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90,
	203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215,
	79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
}

var qrBlocks = []int{
	21, 25, 29, 33, 37, 41, 45, 49,
	53, 57, 61, 65, 69, 73, 77, 81,
	85, 89, 93, 97, 101, 105, 109, 113,
	117, 121, 125, 129, 133, 137, 141, 145,
	149, 153, 157, 161, 165, 169, 173, 177,
}

//QRGenerate генерирует qr
func QRGenerate() {
	fmt.Println(code)

	//Расчет длины массива данных
	length := 0
	for l, i := len(code), 0; i < l; i++ {
		if l > i+1 {
			length += 11
			i++
		} else {
			length += 6
		}
	}

	//Перевод строки в двоичную последовательность
	j := 0
	data := make([]int, length)
	for l, i := len(code), 0; i < l; i++ {
		if l > i+1 {
			mask := 1024
			x := int(ascii[code[i]])*45 + int(ascii[code[i+1]])
			for k := 0; k < 11; k++ {
				if x&mask != 0 {
					data[j] = 1
				}
				mask >>= 1
				j++
			}
			i++
		} else {
			mask := 32
			x := int(ascii[code[i]])
			for k := 0; k < 6; k++ {
				if x&mask != 0 {
					data[j] = 1
				}
				mask >>= 1
				j++
			}
		}
	}

	//Выбор версии QR кода и длины системных данных
	var version int
	var lenSystemData int
	for i := 0; i < 40; i++ {
		max := maxDataH[i]
		if length > max {
			continue
		}
		switch {
		case i < 9:
			if length+13 > max {
				i++
				if i == 9 {
					version = i
					lenSystemData = 15
					break
				}
			}
			version = i
			lenSystemData = 13
		case i >= 9 && i < 26:
			if length+15 > max {
				i++
				if i == 26 {
					version = i
					lenSystemData = 17
					break
				}
			}
			version = i
			lenSystemData = 15
		case i >= 26 && i < 40:
			if length+17 > max {
				i++
				if i == 40 {
					fmt.Println("Data is oversize")
					break
				}
			}
			version = i
			lenSystemData = 17
		}
		break
	}

	//Запись системных данных в начало массива
	newData := make([]int, maxDataH[version])
	newData[2] = 1
	mask := 1 << (lenSystemData - 4 - 1)
	countSymbol := len(code)
	for i := 4; i < lenSystemData; i++ {
		if countSymbol&mask != 0 {
			newData[i] = 1
		}
		mask >>= 1
	}
	copy(newData[lenSystemData:], data)

	//Дозаполнение пустышками до необходимой длины
	minMultiplyLenData := lenSystemData + length
	for minMultiplyLenData%8 != 0 {
		minMultiplyLenData++
	}
	f := true
	for i := minMultiplyLenData; i < maxDataH[version]; {
		switch f {
		case true:
			newData[i] = 1
			newData[i+1] = 1
			newData[i+2] = 1
			newData[i+4] = 1
			newData[i+5] = 1
			i += 8
			f = false
		case false:
			newData[i+3] = 1
			newData[i+7] = 1
			i += 8
			f = true
		}
	}

	//Пстроение блоков
	fullLength := 0
	block := blocksH[version]
	byteData := make([][]int, block)
	maxByte := maxDataH[version] / 8
	size, resid := 0, 0
	if block != 1 {
		k := 0
		size, resid = maxByte/block, maxByte%block
		for i := 0; i < block; i++ {
			if resid >= block-i {
				byteData[i] = make([]int, size+1)
				fullLength += size + 1
			} else {
				byteData[i] = make([]int, size)
				fullLength += size
			}
			for j := 0; j < len(byteData[i]); j++ {
				x := newData[k] * (2 << 6)
				x += newData[k+1] * (2 << 5)
				x += newData[k+2] * (2 << 4)
				x += newData[k+3] * (2 << 3)
				x += newData[k+4] * (2 << 2)
				x += newData[k+5] * (2 << 1)
				x += newData[k+6] * (2 << 0)
				x += newData[k+7] * (1)
				byteData[i][j] = x
				k += 8
			}
		}
	} else {
		count := maxByte
		byteData[0] = make([]int, count)
		fullLength += count
		for i := 0; i < count; i++ {
			x := newData[i*8] * (2 << 6)
			x += newData[i*8+1] * (2 << 5)
			x += newData[i*8+2] * (2 << 4)
			x += newData[i*8+3] * (2 << 3)
			x += newData[i*8+4] * (2 << 2)
			x += newData[i*8+5] * (2 << 1)
			x += newData[i*8+6] * (2 << 0)
			x += newData[i*8+7] * (1)
			byteData[0][i] = x
		}
	}

	//Создание байт коррекции
	countByteCorect := byteCorectH[version]
	polinomCorect := polinom[countByteCorect]
	lenPolinim := len(polinomCorect)
	corectData := make([][]int, block)
	for i := range corectData {
		lenBlock := len(byteData[i])
		if lenBlock > lenPolinim {
			corectData[i] = make([]int, lenBlock)
			fullLength += lenBlock
		} else {
			corectData[i] = make([]int, lenPolinim)
			fullLength += lenPolinim
		}
		copy(corectData[i], byteData[i])
		for range byteData[i] {
			x := corectData[i][0]
			copy(corectData[i], corectData[i][1:])
			if x == 0 {
				continue
			}
			x = reversefieldGalua[x]
			for j := 0; j < countByteCorect; j++ {
				y := polinomCorect[j] + x
				if y > 254 {
					y %= 255
				}
				corectData[i][j] ^= fieldGalua[y]
			}
		}
	}

	//Групирование блоков данных
	i := 0
	data = make([]int, fullLength)
	for j := 0; j < size+1; j++ {
		for _, v := range byteData {
			if len(v) > j {
				data[i] = v[j]
				i++
			}
		}
	}
	for j := 0; true; {
		f := true
		for _, v := range corectData {
			if len(v) > j {
				f = false
				data[i] = v[j]
				i++
			}
		}
		j++
		if f {
			break
		}
	}

	fmt.Println(version)
	fmt.Println(newData)
	fmt.Println(byteData)
	fmt.Println(corectData)
	fmt.Println(data)

	//Вывод изображения
	file, err := os.OpenFile("qrtest.png", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	size = qrBlocks[version] + 4
	rect := image.Rect(0, 0, size, size)
	img := image.NewGray(rect)

	//Рамки
	border(img, size)
	//Поисковые маячки
	searchPoint(img, 4, 4)
	searchPoint(img, 4, size-4-7)
	searchPoint(img, size-4-7, 4)
	//Полоса синханизации
	syncLine(img, 4+6, size-4-6)
	//Информация о маске
	maskInfo(img, 4, 4, size, 5769)

	if err := png.Encode(file, img); err != nil {
		fmt.Println(err)
		return
	}
}

//Рамки
func border(img *image.Gray, size int) {
	for x := 0; x < size; x++ {
		for y := 0; y < 4; y++ {
			img.Set(x, y, color.White)
			img.Set(y, x, color.White)
		}
		for y := size - 4; y < size; y++ {
			img.Set(x, y, color.White)
			img.Set(y, x, color.White)
		}
	}
}

//Поисковые маячки
func searchPoint(img *image.Gray, x, y int) {
	for i := 1; i < 6; i++ {
		img.Set(x+i, y+1, color.White)
		img.Set(y+1, x+i, color.White)
		img.Set(x+i, y+5, color.White)
		img.Set(y+5, x+i, color.White)
	}
	for i := -1; i < 8; i++ {
		img.Set(x+i, y-1, color.White)
		img.Set(y-1, x+i, color.White)
		img.Set(x+i, y+7, color.White)
		img.Set(y+7, x+i, color.White)
	}
}

//Полосы синхранизации
func syncLine(img *image.Gray, start, end int) {
	f := true
	for i := start; i < end; i++ {
		if f {
			f = false
		} else {
			img.Set(start, i, color.White)
			img.Set(i, start, color.White)
			f = true
		}
	}
}

//Информация о версии и маске
func maskInfo(img *image.Gray, x, y, size int, code int) {
	mask := 16384
	for i, f := 0, 0; i < 15; i, f = i+1, f+1 {
		if i > 7 {
			if code&mask == 0 {
				if f == 10 {
					f++
				}
				img.Set(x+8, y+16-f, color.White)
				img.Set(size-x-15+i, y+8, color.White)
			}
			mask >>= 1
		} else {
			if code&mask == 0 {
				if f == 6 {
					f++
				}
				img.Set(x+f, y+8, color.White)
				if i == 7 {
					img.Set(size-x-16+i, y+8, color.White)
				} else {
					img.Set(x+8, size-y-1-i, color.White)
				}
			}
			mask >>= 1
		}
	}
	img.Set(x+8, size-y-8, color.Black)
	//img.Set(x+i, y+8, color.White)
	//img.Set(i, start, color.White)
}
