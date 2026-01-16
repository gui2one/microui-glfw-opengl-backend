#version 460
layout(location=0) in vec2 aPos;
layout(location=1) in vec2 aUvs;
out vec2 fUvs;

void main() {

    gl_Position = vec4(vec3(aPos,0.0), 1.0);
    fUvs = aUvs;
}