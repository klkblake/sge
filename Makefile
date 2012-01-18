include $(GOROOT)/src/Make.inc

TARG=sge
GOFILES=\
	world.go\
	node.go\
	view.go\
	assets.go\
	events.go\
	texture.go\
	shader.go\
	program.go\
	mesh.go\
	buffer.go\
	skybox.go\
	renderer.go\
	mat4.go\
	mat4stack.go\
	error.go\

include $(GOROOT)/src/Make.pkg
