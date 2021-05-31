#include "corto_api.h"
#include "corto.h"

#include <memory>

#ifdef __cplusplus
extern "C" {
#endif

struct _corto_encoder_t {
  std::shared_ptr<crt::Encoder> ptr;
};

corto_encoder_t *corto_new_encoder(uint32_t _nvert, uint32_t _nface,
                                   uint32_t entropy) {
  return new _corto_encoder_t{std::make_shared<crt::Encoder>(
      _nvert, _nface, static_cast<crt::Stream::Entropy>(entropy))};
}

void corto_encoder_free(corto_encoder_t *enc) {
  if (enc)
    delete enc;
}

_Bool corto_encoder_add_positions(corto_encoder_t *enc, const float *buffer,
                                  float q, float *o) {
  if (o) {
    return enc->ptr->addPositions(buffer, q, crt::Point3f(o[0], o[1], o[2]));
  } else {
    return enc->ptr->addPositions(buffer, q);
  }
  return false;
}

_Bool corto_encoder_add_positions_index32(corto_encoder_t *enc,
                                          const float *buffer,
                                          const uint32_t *index, float q,
                                          float *o) {
  if (o) {
    return enc->ptr->addPositions(buffer, index, q,
                                  crt::Point3f(o[0], o[1], o[2]));
  } else {
    return enc->ptr->addPositions(buffer, index, q);
  }
  return false;
}

_Bool corto_encoder_add_positions_index16(corto_encoder_t *enc,
                                          const float *buffer,
                                          const uint16_t *index, float q,
                                          float *o) {
  if (o) {
    return enc->ptr->addPositions(buffer, index, q,
                                  crt::Point3f(o[0], o[1], o[2]));
  } else {
    return enc->ptr->addPositions(buffer, index, q);
  }
  return false;
}

_Bool corto_encoder_add_positions_bits(corto_encoder_t *enc,
                                       const float *buffer, int bits) {
  return enc->ptr->addPositionsBits(buffer, bits);
}

_Bool corto_encoder_add_positions_bits_index32(corto_encoder_t *enc,
                                               const float *buffer,
                                               uint32_t *index, int bits) {
  return enc->ptr->addPositionsBits(buffer, index, bits);
}

_Bool corto_encoder_add_positions_bits_index16(corto_encoder_t *enc,
                                               const float *buffer,
                                               uint16_t *index, int bits) {
  return enc->ptr->addPositionsBits(buffer, index, bits);
}

_Bool corto_encoder_add_normals_float(corto_encoder_t *enc, const float *buffer,
                                      int bits, uint32_t prediction) {
  return enc->ptr->addNormals(
      buffer, bits, static_cast<crt::NormalAttr::Prediction>(prediction));
}

_Bool corto_encoder_add_normals_short(corto_encoder_t *enc,
                                      const int16_t *buffer, int bits,
                                      uint32_t prediction) {
  return enc->ptr->addNormals(
      buffer, bits, static_cast<crt::NormalAttr::Prediction>(prediction));
}

_Bool corto_encoder_add_colors(corto_encoder_t *enc,
                               const unsigned char *buffer, int rbits,
                               int gbits, int bbits, int abits) {
  return enc->ptr->addColors(buffer, rbits, gbits, bbits, abits);
}

_Bool corto_encoder_add_colors3(corto_encoder_t *enc,
                                const unsigned char *buffer, int rbits,
                                int gbits, int bbits) {
  return enc->ptr->addColors3(buffer, rbits, gbits, bbits);
}

_Bool corto_encoder_add_uvs(corto_encoder_t *enc, const float *buffer,
                            float q) {
  return enc->ptr->addUvs(buffer, q);
}

_Bool corto_encoder_add_attribute(corto_encoder_t *enc, const char *name,
                                  const char *buffer, uint32_t format,
                                  int components, float q, uint32_t strategy) {
  return enc->ptr->addAttribute(
      name, buffer, static_cast<crt::VertexAttribute::Format>(format),
      components, q, strategy);
}

void corto_encoder_add_group(corto_encoder_t *enc, int end_triangle) {
  enc->ptr->addGroup(end_triangle);
}

void corto_encoder_add_group_props(corto_encoder_t *enc, int end_triangle,
                                   char **props_keys, char **props_values,
                                   int len) {
  std::map<std::string, std::string> props;
  for (int i = 0; i < len; i++) {
    props.emplace(props_keys[i], props_values[i]);
  }
  enc->ptr->addGroup(end_triangle, props);
}

size_t corto_encoder_encode(corto_encoder_t *enc) { 
  enc->ptr->encode();
  return enc->ptr->stream.size();  
}

void corto_encoder_get_data(corto_encoder_t *enc, char *data, size_t len) {
  memcpy(data, enc->ptr->stream.data(), len);
}

struct _corto_decoder_t {
  std::shared_ptr<crt::Decoder> ptr;
};

corto_decoder_t *corto_new_decoder(int len, const uint8_t *input) {
  return new _corto_decoder_t{std::make_shared<crt::Decoder>(len, input)};
}

void corto_decoder_free(corto_decoder_t *dec) {
  if (dec)
    delete dec;
}

_Bool corto_encoder_has_attr(corto_decoder_t *dec, const char *name) {
  return dec->ptr->hasAttr(name);
}

_Bool corto_encoder_set_positions(corto_decoder_t *dec, float *buffer) {
  return dec->ptr->setPositions(buffer);
}

_Bool corto_encoder_set_normals_float(corto_decoder_t *dec, float *buffer) {
  return dec->ptr->setNormals(buffer);
}

_Bool corto_encoder_set_normals_short(corto_decoder_t *dec, int16_t *buffer) {
  return dec->ptr->setNormals(buffer);
}

_Bool corto_encoder_set_uvs(corto_decoder_t *dec, float *buffer) {
  return dec->ptr->setUvs(buffer);
}

_Bool corto_encoder_set_colors(corto_decoder_t *dec, uint8_t *buffer,
                               int components) {
  return dec->ptr->setColors(buffer, components);
}

_Bool corto_encoder_set_attribute(corto_decoder_t *dec, const char *name,
                                  char *buffer, uint32_t format) {
  return dec->ptr->setAttribute(
      name, buffer, static_cast<crt::VertexAttribute::Format>(format));
}

void corto_encoder_set_index32(corto_decoder_t *dec, uint32_t *buffer) {
  dec->ptr->setIndex(buffer);
}

void corto_encoder_set_index16(corto_decoder_t *dec, uint16_t *buffer) {
  dec->ptr->setIndex(buffer);
}

void corto_encoder_decode(corto_decoder_t *dec) { dec->ptr->decode(); }

#ifdef __cplusplus
}
#endif
