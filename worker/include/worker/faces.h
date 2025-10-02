// Copyright 2025 Omkar Prabhu
#pragma once

// collision between iamgemagick and onnxruntime
#ifdef IsNaN
#undef IsNaN
#endif
#include <onnxruntime_cxx_api.h>

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
  ONNXModel(
      const std::string& detection_model = "models/faces_det/scrfd_2.5g.onnx",
      float detection_threshold = 0.8,
      const std::string& recognition_model =
          "models/faces_rec/webface_r50.onnx")
      : detection_model_(detection_model),
        detection_threshold_(detection_threshold),
        recognition_model_(recognition_model),
        detection_session_(nullptr),
        recognition_session_(nullptr) {
    Ort::Env env(ORT_LOGGING_LEVEL_WARNING, "faces");
    Ort::SessionOptions options;
    options.SetGraphOptimizationLevel(GraphOptimizationLevel::ORT_ENABLE_BASIC);
    detection_session_ = Ort::Session(env, detection_model.c_str(), options);
    recognition_session_ =
        Ort::Session(env, recognition_model.c_str(), options);
  }
  std::vector<std::pair<std::string, std::vector<float>>> Run(
      const std::string& file_path) override;

 private:
  std::string detection_model_;
  float detection_threshold_;
  std::string recognition_model_;
  Ort::Env env_;
  Ort::Session detection_session_;
  Ort::Session recognition_session_;
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
