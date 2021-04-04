package image_filter

import (
	"image"
	"image/color"
)

func Denoize(newImage *image.RGBA, level int) {
	high, low := 10*level, -5*level
	height, width := newImage.Rect.Dy(), newImage.Rect.Dx()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r0, g0, b0, a0 := newImage.At(x, y).RGBA()
			rr0, gg0, bb0, aa0 := int(r0>>8), int(g0>>8), int(b0>>8), int(a0>>8)
			k, k0 := 1, 1
			R, G, B, A := rr0, gg0, bb0, aa0
			R0, G0, B0, A0 := rr0, gg0, bb0, aa0
			for y0 := y - 1; y0 < y+2; y0++ {
				for x0 := x - 1; x0 < x+2; x0++ {
					if x0 < width && y0 < height && x0 > -1 && y0 > -1 {
						r, g, b, a := newImage.At(x0, y0).RGBA()
						rr, gg, bb, aa := int(r>>8), int(g>>8), int(b>>8), int(a>>8)
						R0 += rr
						G0 += gg
						B0 += bb
						A0 += aa
						k0++
						if rd, gd, bd := rr0-rr, gg0-gg, bb0-bb; (rd < high && rd > low) || (gd < high && gd > low) || (bd < high && bd > low) {
							R += rr
							G += gg
							B += bb
							A += aa
							k++
						}
					}
				}
			}
			if k == 2 {
				newImage.Set(x, y, color.RGBA{R: uint8(R0 / k0), G: uint8(G0 / k0), B: uint8(B0 / k0), A: uint8(A0 / k0)})
			} else {
				newImage.Set(x, y, color.RGBA{R: uint8(R / k), G: uint8(G / k), B: uint8(B / k), A: uint8(A / k)})
			}
		}
	}
}
