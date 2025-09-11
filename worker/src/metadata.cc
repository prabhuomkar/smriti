// Copyright 2025 Omkar Prabhu
#include "worker/metadata.h"

#include <spdlog/spdlog.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/components.h"

namespace components {

namespace metadata {

std::unordered_map<std::string, std::string> Metadata::Extract(
    const std::string& user_id, const std::string& mediaitem_id,
    const std::string& file_path) {
  return exif_tool_client_->Extract(file_path);
}

std::shared_ptr<Metadata> Init(const ComponentConfig& config) {
  return std::make_shared<Metadata>(std::make_shared<ExifToolClient>());
}

} // namespace metadata

} // namespace components
