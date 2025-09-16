// Copyright 2025 Omkar Prabhu
#include "worker/preview_thumbnail.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <vector>

#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace previewthumbnail {

std::shared_ptr<PreviewThumbnail> Init(const ComponentConfig& config,
                                       std::shared_ptr<APIClient> api_client) {
  simdjson::ondemand::parser parser;
  simdjson::padded_string padded_config =
      simdjson::padded_string(config.params);
  simdjson::ondemand::document doc = parser.iterate(padded_config);
  int image_quality = static_cast<int>(doc["image_quality"].get_int64());
  int thumbnail_size = static_cast<int>(doc["thumbnail_size"].get_int64());
  int placeholder_size = static_cast<int>(doc["placeholder_size"].get_int64());
  return std::make_shared<PreviewThumbnail>(
      std::make_shared<ImageConverterClient>(), api_client, image_quality,
      thumbnail_size, placeholder_size);
}

std::unordered_map<std::string, std::string> PreviewThumbnail::Generate(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path,
    const std::string& type) {
  std::unordered_map<std::string, std::string> result;

  // default value
  MediaItemStatus status = MediaItemStatus::FAILED;

  MediaItemPreviewThumbnailRequest request;
  request.set_userid(user_id);
  request.set_mediaitemid(mediaitem_id);
  request.set_sourcepath(file_path);

  try {
    if (type == MediaItemType_Name(MediaItemType::PHOTO)) {
      result["preview_url"] = image_converter_client_->Convert(
          file_path, file_path + "-preview", image_quality_, 0);
    } else if (type == MediaItemType_Name(MediaItemType::VIDEO)) {
      result["preview_url"] = "";
    }
    request.set_previewpath(result["preview_url"]);

    result["thumbnail_url"] = image_converter_client_->Convert(
        result["preview_url"], file_path + "-thumbnail", image_quality_,
        thumbnail_size_);
    request.set_thumbnailpath(result["thumbnail_url"]);

    result["placeholder"] = image_converter_client_->Convert(
        result["thumbnail_url"], "", image_quality_, placeholder_size_);
    request.set_placeholder(result["placeholder"]);

    status = MediaItemStatus::READY;
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error extracting preview thumbnail: {}", e.what());
  }

  result["status"] = MediaItemStatus_Name(status);
  request.set_status(status);
  api_client_->SaveMediaItemPreviewThumbnail(request);

  return result;
}

} // namespace previewthumbnail

} // namespace components
