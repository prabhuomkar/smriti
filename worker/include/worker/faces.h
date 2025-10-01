// Copyright 2025 Omkar Prabhu
#pragma once

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

class ModelInferenceInterface {
 public:
  virtual ~ModelInferenceInterface() = default;
  virtual std::vector<std::pair<std::string, std::vector<float>>> Run(
      const std::string& file_path) = 0;
};

class ONNXModel : public ModelInferenceInterface {
 public:
  std::vector<std::pair<std::string, std::vector<float>>> Run(
      const std::string& file_path) override {
    return {};
  }
};

class Faces {
 public:
  Faces() {}
  virtual std::unordered_map<std::string, std::string> Extract(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path);
};

class ONNX : public Faces {
 public:
  explicit ONNX(std::shared_ptr<ModelInferenceInterface> model,
                std::shared_ptr<APIClient> api_client)
      : model_(model), api_client_(api_client) {}
  std::unordered_map<std::string, std::string> Extract(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path) override;

 private:
  std::shared_ptr<ModelInferenceInterface> model_;
  std::shared_ptr<APIClient> api_client_;
};

std::shared_ptr<Faces> Init(const ComponentConfig& config,
                            std::shared_ptr<APIClient> api_client);
} // namespace faces

} // namespace components
