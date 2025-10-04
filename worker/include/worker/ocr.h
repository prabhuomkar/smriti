// Copyright 2025 Omkar Prabhu
#pragma once

#include <spdlog/spdlog.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "src/api/pipelines/ocr.h"
#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace ocr {

class ModelInferenceInterface {
 public:
  virtual ~ModelInferenceInterface() = default;
  virtual std::vector<std::pair<std::string, float>> Run(
      const std::string& file_path) = 0;
};

class PaddlePaddleModel : public ModelInferenceInterface {
 public:
  PaddlePaddleModel(
      bool use_doc_orientation_classify = false,
      const std::string& doc_orientation_classify_model_name =
          "PP-LCNet_x1_0_doc_ori",
      const std::string& doc_orientation_classify_model_dir =
          "ocr_doc_orient/PP-LCNet_x1_0_doc_ori_infer",
      bool use_textline_orientation = false,
      const std::string& textline_orientation_model_name =
          "PP-LCNet_x0_25_textline_ori",
      const std::string& textline_orientation_model_dir =
          "ocr_text_line_orient/PP-LCNet_x0_25_textline_ori_infer",
      bool use_doc_unwarping = false,
      const std::string& doc_unwarping_model_name = "UVDoc",
      const std::string& doc_unwarping_model_dir =
          "ocr_text_unwrap/UVDoc_infer",
      const std::string& text_detection_model_name = "PP-OCRv5_mobile_det",
      const std::string& text_detection_model_dir =
          "ocr_text_det/PP-OCRv5_mobile_det_infer",
      const std::string& text_recognition_model_name = "PP-OCRv5_mobile_rec",
      const std::string& text_recognition_model_dir =
          "ocr_text_rec/PP-OCRv5_mobile_rec_infer")
      : params_(
            {.use_doc_orientation_classify = use_doc_orientation_classify,
             .doc_orientation_classify_model_name =
                 doc_orientation_classify_model_name,
             .doc_orientation_classify_model_dir =
                 doc_orientation_classify_model_dir,
             .use_textline_orientation = use_textline_orientation,
             .textline_orientation_model_name = textline_orientation_model_name,
             .textline_orientation_model_dir = textline_orientation_model_dir,
             .use_doc_unwarping = use_doc_unwarping,
             .doc_unwarping_model_name = doc_unwarping_model_name,
             .doc_unwarping_model_dir = doc_unwarping_model_dir,
             .text_detection_model_name = text_detection_model_name,
             .text_detection_model_dir = text_detection_model_dir,
             .text_recognition_model_name = text_recognition_model_name,
             .text_recognition_model_dir = text_recognition_model_dir}),
        model_(params_) {}
  std::vector<std::pair<std::string, float>> Run(
      const std::string& file_path) override;

 private:
  PaddleOCRParams params_;
  PaddleOCR model_;
};

class OCR {
 public:
  OCR() {}
  virtual std::unordered_map<std::string, std::string> Extract(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path);
};

class PaddlePaddle : public OCR {
 public:
  explicit PaddlePaddle(std::shared_ptr<ModelInferenceInterface> model,
                        std::shared_ptr<APIClient> api_client)
      : model_(model), api_client_(api_client) {}
  std::unordered_map<std::string, std::string> Extract(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path) override;

 private:
  std::shared_ptr<ModelInferenceInterface> model_;
  std::shared_ptr<APIClient> api_client_;
};

std::shared_ptr<OCR> Init(const std::string& models_dir,
                          const ComponentConfig& config,
                          std::shared_ptr<APIClient> api_client);
} // namespace ocr

} // namespace components
