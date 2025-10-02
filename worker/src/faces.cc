// Copyright 2025 Omkar Prabhu
#include "worker/faces.h"

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

namespace faces {

std::vector<std::pair<std::string, std::vector<float>>> ONNXModel::Run(
    const std::string& file_path) {
  return {};
}

std::unordered_map<std::string, std::string> Faces::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  return {};
}

std::unordered_map<std::string, std::string> ONNX::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  if (file_path == "") {
    return {};
  }

  SPDLOG_DEBUG("detecting faces and getting embeddings for file: {}",
               file_path);

  std::vector<std::pair<std::string, std::vector<float>>> response;
  try {
    response = model_->Run(file_path);
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error extracting faces: {}", e.what());
  }

  std::unordered_map<std::string, std::string> result;
  std::string thumbnails = "";

  MediaItemFacesRequest request;
  request.set_userid(user_id);
  request.set_mediaitemid(mediaitem_id);
  for (const auto& [thumbnail, embedding] : response) {
    request.add_thumbnails(thumbnail);
    MediaItemEmbedding* mediaitem_embedding = request.add_embeddings();
    mediaitem_embedding->mutable_embedding()->Add(embedding.begin(),
                                                  embedding.end());
    thumbnails += ((thumbnails == "" ? "" : ",") + thumbnail);
  }
  api_client_->SaveMediaItemFaces(request);

  if (thumbnails != "") {
    result["faces"] = thumbnails;
  }

  return result;
}

std::shared_ptr<Faces> Init(const ComponentConfig& config,
                            std::shared_ptr<APIClient> api_client) {
  if (config.source == "onnx") {
    return std::make_shared<ONNX>(
        std::make_shared<ONNXModel>(
            "../../../models/faces_det/scrfd_2.5g.onnx", 0.8,
            "../../../models/faces_rec/webface_r50.onnx"),
        api_client);
  }
  return nullptr;
}

} // namespace faces

} // namespace components
