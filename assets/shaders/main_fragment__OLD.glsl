#version 460
uniform sampler2D uTexture;
in vec2 fUvs;
in vec3 fColor;
out vec4 frag_color;
void main() {

    vec4 tex = texture(uTexture, fUvs);
    float sharpAlpha = smoothstep(0.2, 0.8, tex.a);
    frag_color = vec4(fColor, sharpAlpha);
}