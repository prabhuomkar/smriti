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
                              int image_quality, int image_size) = 0;
};

class ImageConverterClient : public ImageConverterClientInterface {
 public:
  ImageConverterClient() {}
  ~ImageConverterClient() override {}
  std::string Convert(const std::string& input_file_path,
                      const std::string& output_file_path, int image_quality,
                      int image_size) override;
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
