// Copyright 2025 Omkar Prabhu
#include "worker/ocr.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace ocr {

std::vector<std::pair<std::string, float>> PaddlePaddleModel::Run(
    const std::string& file_path) {
  return {};
}

std::unordered_map<std::string, std::string> OCR::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  return {};
}

std::unordered_map<std::string, std::string> PaddlePaddle::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  if (file_path == "") {
    return {};
  }

  SPDLOG_DEBUG("detecting text from file: {}", file_path);

  std::vector<std::pair<std::string, float>> response;
  try {
    response = model_->Run(file_path);
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error extracting text: {}", e.what());
  }

  std::unordered_map<std::string, std::string> result;

  for (const auto& [text, score] : response) {
    // TODO(omkar): check score threshold
    result["detected_text"] +=
        ((result["detected_text"] == "" ? "" : " ") + text);
  }

  return result;
}

std::shared_ptr<OCR> Init(const std::string& models_dir,
                          const ComponentConfig& config,
                          std::shared_ptr<APIClient> api_client) {
  if (config.source == "paddlepaddle") {
    return std::make_shared<PaddlePaddle>(std::make_shared<PaddlePaddleModel>(),
                                          api_client);
  }
  return nullptr;
}

} // namespace ocr

} // namespace components
