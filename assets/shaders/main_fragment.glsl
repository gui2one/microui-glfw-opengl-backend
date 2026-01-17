#version 460
uniform sampler2D uTexture;
in vec2 fUvs;
in vec3 fColor;
out vec4 frag_color;
void main() {

    vec4 tex = texture(uTexture, fUvs);
    // frag_color = vec4(fUvs.x,fUvs.y,1.0,1.0);
    frag_color = tex * vec4(fColor, 1.0);

}