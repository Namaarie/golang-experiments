//go:build ignore

//kage:unit pixels

package main

var Time float

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {

	uv := vec2(dstPos.xy / imageDstSize())

	return vec4(uv, float((int(Time)%60)/60), 1)
}
