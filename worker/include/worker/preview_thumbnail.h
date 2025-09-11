// Copyright 2025 Omkar Prabhu
#pragma once

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/components.h"

namespace components {

namespace previewthumbnail {

class PreviewThumbnail {
 public:
  explicit PreviewThumbnail(int thumbnail_size = 256, int placeholder_size = 4)
      : thumbnail_size_(thumbnail_size), placeholder_size_(placeholder_size) {}
  std::unordered_map<std::string, std::string> GenerateImage(
      const std::string& user_id, const std::string& mediaitem_id,
      const std::string& file_path, const std::string& type);
  std::unordered_map<std::string, std::string> GenerateVideo(
      const std::string& user_id, const std::string& mediaitem_id,
      const std::string& file_path, const std::string& type);

 private:
  int thumbnail_size_;
  int placeholder_size_;
};

std::shared_ptr<PreviewThumbnail> Init(const ComponentConfig& config);

} // namespace previewthumbnail

} // namespace components
