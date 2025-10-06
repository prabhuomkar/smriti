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

std::string ImageConverterClient::Convert(const std::string& input_file_path,
                                          const std::string& output_file_path,
                                          int image_size) {
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

    libraw_processed_image_t* image = raw_processor.dcraw_make_mem_image(&ret);
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
  magick_img.quality(image_quality_);
  magick_img.magick("JPEG");

  if (output_file_path == "") {
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

std::shared_ptr<PreviewThumbnail> Init(const ComponentConfig& config,
                                       std::shared_ptr<APIClient> api_client) {
  simdjson::ondemand::parser parser;
  simdjson::padded_string padded_config =
      simdjson::padded_string(config.params);
  simdjson::ondemand::document doc = parser.iterate(padded_config);
  int image_quality = 50;
  if (doc["image_quality"].error() == simdjson::SUCCESS) {
    image_quality = static_cast<int>(doc["image_quality"].get_int64());
  }
  int thumbnail_size = 256;
  if (doc["thumbnail_size"].error() == simdjson::SUCCESS) {
    thumbnail_size = static_cast<int>(doc["thumbnail_size"].get_int64());
  }
  int placeholder_size = 2;
  if (doc["placeholder_size"].error() == simdjson::SUCCESS) {
    placeholder_size = static_cast<int>(doc["placeholder_size"].get_int64());
  }
  return std::make_shared<PreviewThumbnail>(
      std::make_shared<ImageConverterClient>(image_quality), api_client,
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
          file_path, file_path + "-preview", 0);
    } else if (type == MediaItemType_Name(MediaItemType::VIDEO)) {
      result["preview_url"] = "";
    }
    request.set_previewpath(result["preview_url"]);

    result["thumbnail_url"] = image_converter_client_->Convert(
        result["preview_url"], file_path + "-thumbnail", thumbnail_size_);
    request.set_thumbnailpath(result["thumbnail_url"]);

    result["placeholder"] = image_converter_client_->Convert(
        result["thumbnail_url"], "", placeholder_size_);
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
