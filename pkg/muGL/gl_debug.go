//go:build gldebug

package muGL

import (
	"log"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

func SetupGLDebug() {
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)

	gl.DebugMessageCallback(func(
		source uint32,
		gltype uint32,
		id uint32,
		severity uint32,
		length int32,
		message string,
		userParam unsafe.Pointer,
	) {
		switch severity {
		case gl.DEBUG_SEVERITY_HIGH:
			log.Printf("[HIGH] %v %s !!!!", gltype, message)
			// panic(message)
		case gl.DEBUG_SEVERITY_MEDIUM:
			log.Printf("[MEDIUM] %v %s !!", gltype, message)
			// panic(message)
		case gl.DEBUG_SEVERITY_LOW:
			log.Printf("[LOW] %s !", message)
			// panic(message)
		case gl.DEBUG_SEVERITY_NOTIFICATION:
			// log.Printf("[NOTIFICATION] %s", message)
			break
		}
	}, nil)
}
