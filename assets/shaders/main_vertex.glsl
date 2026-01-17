#version 460
layout(location=0) in vec2 aPos;
layout(location=1) in vec2 aUvs;
layout(location=2) in vec3 aColor;
out vec2 fUvs;
out vec3 fColor;
uniform mat4 uProj;
void main() {

    gl_Position = uProj * vec4(vec3(aPos,0.0), 1.0);
    fUvs = aUvs;
    fColor = aColor;
}