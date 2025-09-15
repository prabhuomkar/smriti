// Copyright 2025 Omkar Prabhu
#pragma once

#include <Magick++.h>
#include <libraw/libraw.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace previewthumbnail {

inline std::string Base64Encode(const std::vector<unsigned char>& data) {
  static const char table[] =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
      "abcdefghijklmnopqrstuvwxyz"
      "0123456789+/";

  std::string result;
  size_t i = 0;
  unsigned int val = 0;
  int valb = -6;

  for (unsigned char c : data) {
    val = (val << 8) + c;
    valb += 8;
    while (valb >= 0) {
      result.push_back(table[(val >> valb) & 0x3F]);
      valb -= 6;
    }
  }

  if (valb > -6) {
    result.push_back(table[((val << 8) >> (valb + 8)) & 0x3F]);
  }

  while (result.size() % 4) {
    result.push_back('=');
  }

  return result;
}

class ImageConverterClientInterface {
 public:
  virtual ~ImageConverterClientInterface() = default;
  virtual std::string Convert(const std::string& input_file_path,
                              const std::string& output_file_path,
                              int image_quality, int image_size,
                              bool encode = false) = 0;
};

class ImageConverterClient : public ImageConverterClientInterface {
 public:
  ImageConverterClient() {}

  ~ImageConverterClient() override {}

  std::string Convert(const std::string& input_file_path,
                      const std::string& output_file_path, int image_quality,
                      int image_size, bool encode) override {
    LibRaw raw_processor;
    Magick::Image magick_img;
    Magick::Blob magick_blob;

    int ret = raw_processor.open_file(input_file_path.c_str());
    if (ret == LIBRAW_SUCCESS) {
      raw_processor.imgdata.params.no_auto_bright = 1;
      raw_processor.imgdata.params.use_camera_wb = 1;
      raw_processor.imgdata.params.use_camera_matrix = 1;
      raw_processor.imgdata.params.output_color = 1;
      if (raw_processor.unpack() != LIBRAW_SUCCESS ||
          raw_processor.dcraw_process() != LIBRAW_SUCCESS) {
        throw std::runtime_error("unpack or dcraw_process failed");
      }

      libraw_processed_image_t* image =
          raw_processor.dcraw_make_mem_image(&ret);
      if (!image) {
        throw std::runtime_error("dcraw_make_mem_image failed");
      }

      SPDLOG_DEBUG("file_path {}", input_file_path);

      magick_blob = Magick::Blob(image->data, image->data_size);
      magick_img = Magick::Image(
          magick_blob, Magick::Geometry(image->width, image->height), 8, "RGB");

      LibRaw::dcraw_clear_mem(image);
      raw_processor.recycle();
    } else {
      magick_img.read(input_file_path);
    }

    int resize_width = magick_img.columns(), resize_height = magick_img.rows();
    if (image_size > 0) {
      if (resize_width > resize_height) {
        resize_height =
            static_cast<int>((image_size * resize_height) / resize_width);
        resize_width = image_size;
      } else {
        resize_width =
            static_cast<int>((image_size * resize_width) / resize_height);
        resize_height = image_size;
      }
    }

    SPDLOG_DEBUG("resizing image to {} {}", resize_width, resize_height);

    magick_img.resize(Magick::Geometry(resize_width, resize_height));
    magick_img.quality(image_quality);
    magick_img.magick("JPEG");

    if (encode) {
      magick_img.write(&magick_blob, "JPEG");
      const unsigned char* magick_blob_data_ptr =
          static_cast<const unsigned char*>(magick_blob.data());
      std::vector<unsigned char> buffer(
          magick_blob_data_ptr, magick_blob_data_ptr + magick_blob.length());
      return Base64Encode(buffer);
    }

    magick_img.write(output_file_path);
    return output_file_path;
  }
};

class PreviewThumbnail {
 public:
  explicit PreviewThumbnail(
      std::shared_ptr<ImageConverterClientInterface> image_converter_client,
      std::shared_ptr<APIClient> api_client, int image_quality = 50,
      int thumbnail_size = 256, int placeholder_size = 2)
      : image_converter_client_(image_converter_client),
        api_client_(api_client),
        image_quality_(image_quality),
        thumbnail_size_(thumbnail_size),
        placeholder_size_(placeholder_size) {}
  std::unordered_map<std::string, std::string> Generate(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path,
      const std::string& type);

 private:
  std::shared_ptr<ImageConverterClientInterface> image_converter_client_;
  std::shared_ptr<APIClient> api_client_;
  int image_quality_;
  int thumbnail_size_;
  int placeholder_size_;
};

std::shared_ptr<PreviewThumbnail> Init(const ComponentConfig& config,
                                       std::shared_ptr<APIClient> api_client);

} // namespace previewthumbnail

} // namespace components
