// Copyright 2025 Omkar Prabhu
#include "worker/preview_thumbnail.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>

#include "worker/components.h"

namespace components {

namespace previewthumbnail {

std::shared_ptr<PreviewThumbnail> Init(const ComponentConfig& config) {
  simdjson::ondemand::parser parser;
  simdjson::padded_string padded_config =
      simdjson::padded_string(config.params);
  simdjson::ondemand::document doc = parser.iterate(padded_config);
  int thumbnail_size = static_cast<int>(doc["thumbnail_size"].get_int64());
  int placeholder_size = static_cast<int>(doc["placeholder_size"].get_int64());
  return std::make_shared<PreviewThumbnail>(thumbnail_size, placeholder_size);
}

std::unordered_map<std::string, std::string> GenerateImage(
    const std::string& id, const std::string& mediaitem_id,
    const std::string& file_path, const std::string& type) {
  std::unordered_map<std::string, std::string> result;
  return result;
}

std::unordered_map<std::string, std::string> GenerateVideo(
    const std::string& id, const std::string& mediaitem_id,
    const std::string& file_path, const std::string& type) {
  std::unordered_map<std::string, std::string> result;
  return result;
}

} // namespace previewthumbnail

} // namespace components
