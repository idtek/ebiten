// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opengl

import (
	"math"

	"github.com/hajimehoshi/ebiten/internal/graphics"
)

var (
	zeroPlus = math.Nextafter32(0, 1)
	oneMinus = math.Nextafter32(1, 0)
)

var (
	vertexShader       shaderType
	fragmentShader     shaderType
	arrayBuffer        bufferType
	elementArrayBuffer bufferType
	dynamicDraw        bufferUsage
	staticDraw         bufferUsage
	Short              DataType
	Float              DataType

	zero             operation
	one              operation
	srcAlpha         operation
	dstAlpha         operation
	oneMinusSrcAlpha operation
	oneMinusDstAlpha operation
)

func convertOperation(op graphics.Operation) operation {
	switch op {
	case graphics.Zero:
		return zero
	case graphics.One:
		return one
	case graphics.SrcAlpha:
		return srcAlpha
	case graphics.DstAlpha:
		return dstAlpha
	case graphics.OneMinusSrcAlpha:
		return oneMinusSrcAlpha
	case graphics.OneMinusDstAlpha:
		return oneMinusDstAlpha
	default:
		panic("not reached")
	}
}

type Context struct {
	locationCache      *locationCache
	screenFramebuffer  Framebuffer // This might not be the default frame buffer '0' (e.g. iOS).
	lastFramebuffer    Framebuffer
	lastTexture        Texture
	lastViewportWidth  int
	lastViewportHeight int
	lastCompositeMode  graphics.CompositeMode
	maxTextureSize     int
	context
}

var theContext *Context

func GetContext() *Context {
	return theContext
}

func (c *Context) BindTexture(t Texture) {
	if c.lastTexture == t {
		return
	}
	c.bindTextureImpl(t)
	c.lastTexture = t
}

func (c *Context) bindFramebuffer(f Framebuffer) {
	if c.lastFramebuffer == f {
		return
	}
	c.bindFramebufferImpl(f)
	c.lastFramebuffer = f
}

func (c *Context) SetViewport(f Framebuffer, width, height int) {
	c.bindFramebuffer(f)
	if c.lastViewportWidth != width || c.lastViewportHeight != height {
		c.setViewportImpl(width, height)
		// glViewport must be called at least at every frame on iOS.
		// As the screen framebuffer is the last render target, next SetViewport should be
		// the first call at a frame.
		if f == c.screenFramebuffer {
			c.lastViewportWidth = 0
			c.lastViewportHeight = 0
		} else {
			c.lastViewportWidth = width
			c.lastViewportHeight = height
		}
	}
}

func (c *Context) ScreenFramebuffer() Framebuffer {
	return c.screenFramebuffer
}

func (c *Context) MaxTextureSize() int {
	if c.maxTextureSize == 0 {
		c.maxTextureSize = c.maxTextureSizeImpl()
	}
	return c.maxTextureSize
}
