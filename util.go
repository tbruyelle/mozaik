package main

import "golang.org/x/mobile/exp/f32"

func identity() *f32.Mat4 {
	id := &f32.Mat4{}
	id.Identity()
	return id
}

func translate(Tx, Ty, Tz float32) *f32.Mat4 {
	//return &f32.Mat4{
	//	{1, 0, 0, Tx},
	//	{0, 1, 0, Ty},
	//	{0, 0, 1, Tz},
	//	{0, 0, 0, 1},
	//}
	ret := &f32.Mat4{}
	ret.Translate(identity(), Tx, Ty, Tz)
	return ret
}

func rotate(angle float32) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Rotate(identity(), f32.Radian(angle), &f32.Vec3{0, 0, 1})
	return ret
}

func scale(scale float32) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Scale(identity(), scale, scale, 0)
	return ret
}

func mul(m1 *f32.Mat4, m2 *f32.Mat4) *f32.Mat4 {
	ret := &f32.Mat4{}
	ret.Mul(m1, m2)
	return ret
}

func ortho(w, h float32) *f32.Mat4 {
	return &f32.Mat4{
		{2. / w, 0, 0, -1},
		{0, -2. / h, 0, 1},
		{0, 0, -1, 0},
		{0, 0, 0, 1},
	}
	//( 2.0/768.0, 0.0, 0.0, -1.0,
	//  0.0, 2.0/1024.0, 0.0, -1.0,
	//  0.0, 0.0, -1.0, 0.0,
	//  0.0, 0.0, 0.0, 1.0);  )
}
