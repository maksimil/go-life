package main

import (
	"runtime"
	"time"
	"unsafe"

	gl "github.com/go-gl/gl/v3.2-core/gl"
	glfw "github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 750
	height = 750
)

const (
	tw = 1000
	th = 1000
)

func main() {
	// initialization
	runtime.LockOSThread()

	window := initGLFW()
	defer glfw.Terminate()

	initOpenGL()

	// creating shader program
	prog := func() uint32 {

		vertexshader, err := compileShader(TRIANGLE_VERTEX, gl.VERTEX_SHADER)
		if err != nil {
			panic(err)
		}

		fragmentshader, err := compileShader(TRIANGLE_FRAGMENT, gl.FRAGMENT_SHADER)
		if err != nil {
			panic(err)
		}

		prog := gl.CreateProgram()
		gl.AttachShader(prog, vertexshader)
		gl.AttachShader(prog, fragmentshader)
		gl.LinkProgram(prog)
		return prog
	}()

	// creating idicies for rendering
	idxs := func() []float32 {
		idxs := []float32{}

		for j := 0; j < th; j++ {
			for i := 0; i < tw; i++ {
				idxs = append(idxs,
					float32(i+j*tw), float32(i+j*(tw+1)),
					float32(i+j*tw), float32(i+j*(tw+1)+1),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)+1),

					float32(i+j*tw), float32(i+j*(tw+1)),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)),
					float32(i+j*tw), float32(i+(j+1)*(tw+1)+1),
				)
			}
		}
		return idxs
	}()

	vbo := mkvbo(idxs)

	vao := mkvao()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	vertAttribPtr(prog, "idx\x00", 1, gl.FLOAT, 8, 4)
	vertAttribPtr(prog, "cell\x00", 1, gl.FLOAT, 8, 0)

	// settting tilesize uniform
	gl.UseProgram(prog)
	tsloc := gl.GetUniformLocation(prog, gl.Str("tilesize\x00"))
	gl.Uniform2ui(tsloc, uint32(tw), uint32(th))

	// creating the texture
	// data initialization
	state := newSwitch([2]int{tw, th})

	state.states[502_510] = 1
	state.states[502_511] = 1
	state.states[502_514] = 1
	state.states[502_515] = 1
	state.states[502_516] = 1

	state.states[501_513] = 1

	state.states[500_511] = 1

	// state texture initialization
	var texture uint32
	gl.GenTextures(1, &texture)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(tw), int32(th), 0,
		gl.RED, gl.UNSIGNED_BYTE, state.getTexData())

	// state texture bind
	gl.UseProgram(prog)
	stateloc := gl.GetUniformLocation(prog, gl.Str("state\x00"))
	gl.Uniform1i(stateloc, 0)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	printerr()
	LOOPTIME := int64(1 * 1000000)
	// redering loop
	for !window.ShouldClose() {
		s := time.Now()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog)

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(idxs)))
		glfw.PollEvents()
		window.SwapBuffers()

		state.update()
		gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(tw), int32(th), 0,
			gl.RED, gl.UNSIGNED_BYTE, state.getTexData())
		elapsed := time.Since(s).Nanoseconds()
		time.Sleep(time.Duration(LOOPTIME - elapsed))
	}
	printerr()
}

type StateSwitch struct {
	states   []uint8
	size     [2]int
	stateidx int
}

func newSwitch(size [2]int) StateSwitch {
	return StateSwitch{
		make([]uint8, size[0]*size[1]*2),
		size, 0,
	}
}

func (state *StateSwitch) get(k int, i int, j int) *uint8 {
	return &state.states[i+state.size[0]*j+state.size[0]*state.size[1]*k]
}

func (state *StateSwitch) getcurr(i int, j int) *uint8 {
	return state.get(state.stateidx, i, j)
}

func (state *StateSwitch) getnext(i int, j int) *uint8 {
	return state.get(1-state.stateidx, i, j)
}

var OFFSETS = [8][2]int{
	{-1, 1}, {0, 1}, {1, 1},
	{-1, 0}, {1, 0},
	{-1, -1}, {0, -1}, {1, -1},
}

func (state *StateSwitch) update() {
	for i := 0; i < state.size[0]; i++ {
		for j := 0; j < state.size[1]; j++ {
			ncount := uint8(0)
			for k := 0; k < 8; k++ {
				in := (i + OFFSETS[k][0] + state.size[0]) % state.size[0]
				jn := (j + OFFSETS[k][1] + state.size[1]) % state.size[1]
				ncount += *state.getcurr(in, jn)
			}

			if (*state.getcurr(i, j)+ncount == 3) || ncount == 3 {
				*state.getnext(i, j) = 1
			} else {
				*state.getnext(i, j) = 0
			}
		}
	}
	state.stateidx = 1 - state.stateidx
}

func (state *StateSwitch) getTexData() unsafe.Pointer {
	return gl.Ptr(&state.states[state.stateidx*state.size[0]*state.size[1]])
}
