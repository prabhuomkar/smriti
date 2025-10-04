// Copyright 2025 Omkar Prabhu
#include "worker/ocr.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <filesystem> // NOLINT
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
  std::vector<std::pair<std::string, float>> result;

  std::filesystem::rename(file_path, file_path + ".jpg");
  auto outputs = model_.Predict(file_path + ".jpg");
  std::filesystem::rename(file_path + ".jpg", file_path);

  const std::string ocr_output_dir = file_path + "_ocr";
  const std::string file_name = std::filesystem::path(file_path).stem();
  const std::string ocr_output_json_file_name =
      ocr_output_dir + "/" + file_name + "_res.json";
  outputs[0]->SaveToJson(ocr_output_dir);

  simdjson::ondemand::parser parser;
  simdjson::padded_string ocr_output_json =
      simdjson::padded_string::load(ocr_output_json_file_name);
  simdjson::ondemand::document doc = parser.iterate(ocr_output_json);

  if (doc["rec_texts"].get_array().error() != simdjson::SUCCESS ||
      doc["rec_scores"].get_array().error() != simdjson::SUCCESS) {
    return {};
  }

  std::vector<std::string> texts;
  std::vector<float> scores;
  for (auto rec_text : doc["rec_texts"].get_array()) {
    texts.emplace_back(std::string(rec_text.get_string().value()));
  }
  for (auto rec_score : doc["rec_scores"].get_array()) {
    scores.emplace_back(static_cast<float>(rec_score.get_double().value()));
  }
  for (size_t i = 0; i < texts.size(); i++) {
    result.push_back({texts[i], scores[i]});
  }

  std::uintmax_t count = std::filesystem::remove_all(ocr_output_dir);
  SPDLOG_DEBUG("removed ocr output directory {} and file count: {}",
               ocr_output_dir, count);

  return result;
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
    result["detected_text"] +=
        ((result["detected_text"] == "" ? "" : " ") + text);
  }

  return result;
}

std::shared_ptr<OCR> Init(const std::string& models_dir,
                          const ComponentConfig& config,
                          std::shared_ptr<APIClient> api_client) {
  if (config.source == "paddlepaddle") {
    simdjson::ondemand::parser parser;
    simdjson::padded_string padded_config =
        simdjson::padded_string(config.params);
    simdjson::ondemand::document doc = parser.iterate(padded_config);
    bool use_doc_orientation_classify = false;
    std::string doc_orientation_classify_model_name = "PP-LCNet_x1_0_doc_ori";
    std::string doc_orientation_classify_model_dir =
        "ocr_doc_orient/PP-LCNet_x1_0_doc_ori_infer";
    bool use_textline_orientation = false;
    std::string textline_orientation_model_name = "PP-LCNet_x0_25_textline_ori";
    std::string textline_orientation_model_dir =
        "ocr_text_line_orient/PP-LCNet_x0_25_textline_ori_infer";
    bool use_doc_unwarping = false;
    std::string doc_unwarping_model_name = "UVDoc";
    std::string doc_unwarping_model_dir = "ocr_text_unwrap/UVDoc_infer";
    std::string text_detection_model_name = "PP-OCRv5_mobile_det";
    std::string text_detection_model_dir =
        "ocr_text_det/PP-OCRv5_mobile_det_infer";
    std::string text_recognition_model_name = "PP-OCRv5_mobile_rec";
    std::string text_recognition_model_dir =
        "ocr_text_rec/PP-OCRv5_mobile_rec_infer";
    std::string detection_model = "faces_det/scrfd_2.5g.onnx";
    float text_rec_score_thresh = 0.9;
    if (doc["use_doc_orientation_classify"].error() == simdjson::SUCCESS) {
      use_doc_orientation_classify = static_cast<bool>(
          doc["use_doc_orientation_classify"].get_bool().value());
    }
    if (doc["doc_orientation_classify_model_name"].error() ==
        simdjson::SUCCESS) {
      doc_orientation_classify_model_name = std::string(
          doc["doc_orientation_classify_model_name"].get_string().value());
    }
    if (doc["doc_orientation_classify_model_dir"].error() ==
        simdjson::SUCCESS) {
      doc_orientation_classify_model_dir = std::string(
          doc["doc_orientation_classify_model_dir"].get_string().value());
    }
    if (doc["use_textline_orientation"].error() == simdjson::SUCCESS) {
      use_textline_orientation =
          static_cast<bool>(doc["use_textline_orientation"].get_bool().value());
    }
    if (doc["textline_orientation_model_name"].error() == simdjson::SUCCESS) {
      textline_orientation_model_name = std::string(
          doc["textline_orientation_model_name"].get_string().value());
    }
    if (doc["textline_orientation_model_dir"].error() == simdjson::SUCCESS) {
      textline_orientation_model_dir = std::string(
          doc["textline_orientation_model_dir"].get_string().value());
    }
    if (doc["use_doc_unwarping"].error() == simdjson::SUCCESS) {
      use_doc_unwarping =
          static_cast<bool>(doc["use_doc_unwarping"].get_bool().value());
    }
    if (doc["doc_unwarping_model_name"].error() == simdjson::SUCCESS) {
      doc_unwarping_model_name =
          std::string(doc["doc_unwarping_model_name"].get_string().value());
    }
    if (doc["doc_unwarping_model_dir"].error() == simdjson::SUCCESS) {
      doc_unwarping_model_dir =
          std::string(doc["doc_unwarping_model_dir"].get_string().value());
    }
    if (doc["text_detection_model_name"].error() == simdjson::SUCCESS) {
      text_detection_model_name =
          std::string(doc["text_detection_model_name"].get_string().value());
    }
    if (doc["text_detection_model_dir"].error() == simdjson::SUCCESS) {
      text_detection_model_dir =
          std::string(doc["text_detection_model_dir"].get_string().value());
    }
    if (doc["text_recognition_model_name"].error() == simdjson::SUCCESS) {
      text_recognition_model_name =
          std::string(doc["text_recognition_model_name"].get_string().value());
    }
    if (doc["text_recognition_model_dir"].error() == simdjson::SUCCESS) {
      text_recognition_model_dir =
          std::string(doc["text_recognition_model_dir"].get_string().value());
    }
    if (doc["text_rec_score_thresh"].error() == simdjson::SUCCESS) {
      text_rec_score_thresh =
          static_cast<float>(doc["text_rec_score_thresh"].get_double().value());
    }

    return std::make_shared<PaddlePaddle>(
        std::make_shared<PaddlePaddleModel>(
            use_doc_orientation_classify, doc_orientation_classify_model_name,
            models_dir + "/" + doc_orientation_classify_model_dir,
            use_textline_orientation, textline_orientation_model_name,
            models_dir + "/" + textline_orientation_model_dir,
            use_doc_unwarping, doc_unwarping_model_name,
            models_dir + "/" + doc_unwarping_model_dir,
            text_detection_model_name,
            models_dir + "/" + text_detection_model_dir,
            text_recognition_model_name,
            models_dir + "/" + text_recognition_model_dir,
            text_rec_score_thresh),
        api_client);
  }

  return nullptr;
}

} // namespace ocr

} // namespace components
