// Copyright 2025 Omkar Prabhu
#include "worker/faces.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <algorithm>
#include <fstream>
#include <memory>
#include <numeric>
#include <optional>
#include <string>
#include <tuple>
#include <unordered_map>
#include <utility>
#include <vector>

#include "opencv2/dnn.hpp"
#include "opencv2/opencv.hpp"
#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace faces {

std::unordered_map<std::string, std::unordered_map<std::string, std::string>>
FAISSClustering::Assign(
    const std::vector<std::tuple<std::string, std::string, std::string,
                                 std::vector<float>>>& input) {
  std::unordered_map<std::string, std::unordered_map<std::string, std::string>>
      result;

  return result;
}

std::vector<std::string> ONNXModel::Detect(const std::string& file_path) {
  std::vector<std::string> result;

  // preprocess
  cv::Mat img = cv::imread(file_path);
  float im_ratio = static_cast<float>(img.rows) / img.cols;
  float model_ratio = 1.0;
  int new_h = (im_ratio > model_ratio) ? 640 : static_cast<int>(640 * im_ratio);
  int new_w = (im_ratio > model_ratio) ? static_cast<int>(640 / im_ratio) : 640;
  float det_scale = static_cast<float>(new_h) / img.rows;
  cv::Mat resized;
  cv::resize(img, resized, cv::Size(new_w, new_h));
  cv::Mat det_img(640, 640, CV_8UC3, cv::Scalar::all(0));
  resized.copyTo(det_img(cv::Rect(0, 0, new_w, new_h)));
  cv::Mat blob =
      cv::dnn::blobFromImage(det_img, 1.0 / 128.0, cv::Size(640, 640),
                             cv::Scalar(127.5, 127.5, 127.5), true, false);
  Ort::AllocatorWithDefaultOptions allocator;
  auto input_name_ptr = detection_session_.GetInputNameAllocated(0, allocator);
  const char* input_name = input_name_ptr.get();
  std::vector<int64_t> input_shape = {1, 3, blob.size[2], blob.size[3]};
  std::vector<float> input_tensor_values(blob.total() * blob.channels());
  std::memcpy(input_tensor_values.data(), blob.data,
              input_tensor_values.size() * sizeof(float));
  Ort::MemoryInfo memory_info =
      Ort::MemoryInfo::CreateCpu(OrtArenaAllocator, OrtMemTypeDefault);
  Ort::Value input_tensor = Ort::Value::CreateTensor<float>(
      memory_info, input_tensor_values.data(), input_tensor_values.size(),
      input_shape.data(), input_shape.size());
  std::vector<const char*> input_names{input_name};
  std::vector<Ort::AllocatedStringPtr> output_name_ptrs;
  std::vector<const char*> output_names;
  for (size_t i = 0; i < detection_session_.GetOutputCount(); ++i) {
    auto output_name_ptr =
        detection_session_.GetOutputNameAllocated(i, allocator);
    output_names.push_back(output_name_ptr.get());
    output_name_ptrs.push_back(std::move(output_name_ptr));
  }

  // inference run
  auto outputs = detection_session_.Run(
      Ort::RunOptions{nullptr}, input_names.data(), &input_tensor, 1,
      output_names.data(), output_names.size());

  // postprocess
  auto info = outputs[0].GetTensorTypeAndShapeInfo();
  std::vector<cv::Rect> all_boxes;
  std::vector<float> all_scores;
  std::vector<int> strides = {8, 16, 32};
  for (size_t stride_idx = 0; stride_idx < strides.size(); stride_idx++) {
    int stride = strides[stride_idx];
    float* scores_ptr = outputs[stride_idx].GetTensorMutableData<float>();
    int num_scores =
        outputs[stride_idx].GetTensorTypeAndShapeInfo().GetElementCount();
    std::vector<float> scores(scores_ptr, scores_ptr + num_scores);
    float* bboxes_ptr = outputs[stride_idx + 3].GetTensorMutableData<float>();
    int num_bboxes =
        outputs[stride_idx + 3].GetTensorTypeAndShapeInfo().GetElementCount();
    std::vector<float> bboxes(bboxes_ptr, bboxes_ptr + num_bboxes);
    for (size_t i = 0; i < bboxes.size(); ++i) {
      bboxes[i] *= stride;
    }
    int h = blob.size[2] / stride, w = blob.size[3] / stride;
    std::vector<std::pair<float, float>> anchors;
    anchors.reserve(h * w * 2); // Reserve for h*w*2 anchors
    for (size_t y = 0; y < h; ++y) {
      for (size_t x = 0; x < w; ++x) {
        float anchor_x = x * stride;
        float anchor_y = y * stride;
        anchors.emplace_back(anchor_x, anchor_y);
        anchors.emplace_back(anchor_x, anchor_y); // Repeat each anchor
      }
    }
    for (size_t i = 0; i < scores.size(); ++i) {
      if (scores[i] >= 0.5f) {
        float x1 = (anchors[i].first - bboxes[i * 4 + 0]) / det_scale;
        float y1 = (anchors[i].second - bboxes[i * 4 + 1]) / det_scale;
        float x2 = (anchors[i].first + bboxes[i * 4 + 2]) / det_scale;
        float y2 = (anchors[i].second + bboxes[i * 4 + 3]) / det_scale;
        all_boxes.emplace_back(
            cv::Rect(cv::Point(static_cast<int>(x1), static_cast<int>(y1)),
                     cv::Point(static_cast<int>(x2), static_cast<int>(y2))));
        all_scores.push_back(scores[i]);
      }
    }
  }
  std::vector<float> x1(all_boxes.size()), y1(all_boxes.size()),
      x2(all_boxes.size()), y2(all_boxes.size());
  std::vector<float> areas(all_boxes.size());
  for (size_t i = 0; i < all_boxes.size(); ++i) {
    x1[i] = static_cast<float>(all_boxes[i].x);
    y1[i] = static_cast<float>(all_boxes[i].y);
    x2[i] = static_cast<float>(all_boxes[i].x + all_boxes[i].width);
    y2[i] = static_cast<float>(all_boxes[i].y + all_boxes[i].height);
    areas[i] = (x2[i] - x1[i] + 1) * (y2[i] - y1[i] + 1);
  }
  std::vector<int> order(all_boxes.size());
  std::iota(order.begin(), order.end(), 0);
  std::sort(order.begin(), order.end(), [&all_scores](int i, int j) {
    return all_scores[i] > all_scores[j];
  });
  std::vector<int> keep;
  while (!order.empty()) {
    int i = order[0];
    keep.push_back(i);
    std::vector<int> new_order;
    for (size_t k = 1; k < order.size(); ++k) {
      int j = order[k];
      float xx1 = std::max(x1[i], x1[j]);
      float yy1 = std::max(y1[i], y1[j]);
      float xx2 = std::min(x2[i], x2[j]);
      float yy2 = std::min(y2[i], y2[j]);
      float w = std::max(0.0f, xx2 - xx1 + 1);
      float h = std::max(0.0f, yy2 - yy1 + 1);
      float inter = w * h;
      float ovr = inter / (areas[i] + areas[j] - inter);
      if (ovr <= 0.4f) {
        new_order.push_back(j);
      }
    }
    order = std::move(new_order);
  }

  // faces from box grids
  cv::Mat original_img = cv::imread(file_path);
  for (int idx : keep) {
    if (all_scores[idx] >= detection_threshold_) {
      cv::Rect box =
          all_boxes[idx] & cv::Rect(0, 0, original_img.cols, original_img.rows);
      if (box.width > 0 && box.height > 0) {
        cv::Mat face = original_img(box);
        std::string face_out_path = file_path + "_face_" + std::to_string(idx);
        std::vector<uchar> face_out_buf;
        std::vector<int> params = {cv::IMWRITE_JPEG_QUALITY, 50};
        if (!cv::imencode(".jpg", face, face_out_buf, params)) {
          SPDLOG_ERROR("error writing mediaitem face thumbnail");
          continue;
        }
        std::ofstream face_out_file(face_out_path, std::ios::binary);
        face_out_file.write(reinterpret_cast<const char*>(face_out_buf.data()),
                            face_out_buf.size());
        face_out_file.close();
        result.push_back(face_out_path);
      }
    }
  }

  return result;
}

std::vector<std::vector<float>> ONNXModel::Recognize(
    const std::vector<std::string>& face_file_paths) {
  std::vector<std::vector<float>> result;

  for (const std::string& face_file_path : face_file_paths) {
    // preprocess
    cv::Mat img = cv::imread(face_file_path);
    cv::resize(img, img, cv::Size(112, 112));
    cv::cvtColor(img, img, cv::COLOR_BGR2RGB);
    img.convertTo(img, CV_32F, 1.0 / 255.0);
    img = (img - 0.5) / 0.5;
    std::vector<cv::Mat> channels(3);
    cv::split(img, channels);
    cv::Mat blob;
    cv::merge(channels, blob);
    blob = blob.reshape(1, 1);
    Ort::AllocatorWithDefaultOptions allocator;
    auto input_name_ptr =
        recognition_session_.GetInputNameAllocated(0, allocator);
    const char* input_name = input_name_ptr.get();
    std::vector<int64_t> input_shape = {1, 3, 112, 112};
    std::vector<float> input_tensor_values(1 * 3 * 112 * 112);
    for (size_t c = 0; c < 3; ++c) {
      for (size_t h = 0; h < 112; ++h) {
        for (size_t w = 0; w < 112; ++w) {
          input_tensor_values[c * 112 * 112 + h * 112 + w] =
              channels[c].at<float>(h, w);
        }
      }
    }
    Ort::MemoryInfo memory_info =
        Ort::MemoryInfo::CreateCpu(OrtArenaAllocator, OrtMemTypeDefault);
    Ort::Value input_tensor = Ort::Value::CreateTensor<float>(
        memory_info, input_tensor_values.data(), input_tensor_values.size(),
        input_shape.data(), input_shape.size());
    std::vector<const char*> input_names{input_name};
    std::vector<Ort::AllocatedStringPtr> output_name_ptrs;
    std::vector<const char*> output_names;
    for (size_t i = 0; i < recognition_session_.GetOutputCount(); ++i) {
      auto output_name_ptr =
          recognition_session_.GetOutputNameAllocated(i, allocator);
      output_names.push_back(output_name_ptr.get());
      output_name_ptrs.push_back(std::move(output_name_ptr));
    }

    // inference run
    auto outputs = recognition_session_.Run(
        Ort::RunOptions{nullptr}, input_names.data(), &input_tensor, 1,
        output_names.data(), output_names.size());

    // postprocess
    float* output_ptr = outputs[0].GetTensorMutableData<float>();
    int output_size = outputs[0].GetTensorTypeAndShapeInfo().GetElementCount();
    std::vector<float> embedding(output_ptr, output_ptr + output_size);
    result.push_back(std::move(embedding));
  }

  return result;
}

std::vector<std::pair<std::string, std::vector<float>>> ONNXModel::Run(
    const std::string& file_path) {
  std::vector<std::string> face_file_paths;
  try {
    face_file_paths = Detect(file_path);
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error detecting faces: {}", e.what());
  }

  std::vector<std::vector<float>> embeddings;
  try {
    embeddings = Recognize(face_file_paths);
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error recognizing faces: {}", e.what());
  }

  std::vector<std::pair<std::string, std::vector<float>>> result;
  result.reserve(face_file_paths.size());

  for (size_t i = 0; i < face_file_paths.size(); ++i) {
    result.emplace_back(std::move(face_file_paths[i]),
                        std::move(embeddings[i]));
  }

  return result;
}

std::unordered_map<std::string, std::string> Faces::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  return {};
}

void Faces::Cluster() {}

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

void ONNX::Cluster() {
  SPDLOG_INFO("getting users for clustering faces");
  UsersResponse users_response = api_client_->GetUsers();
  for (const std::string& user_id : users_response.users()) {
    SPDLOG_DEBUG("getting mediaitem face embeddings for user: {}", user_id);
    MediaItemFaceEmbeddingsRequest req;
    req.set_userid(user_id);
    MediaItemFaceEmbeddingsResponse mfe_response =
        api_client_->GetMediaItemFaceEmbeddings(req);
    std::vector<std::tuple<std::string, std::string, std::string>>
        face_mediaitem_people_id;
    std::vector<std::vector<float>> embeddings;
    for (auto mfe : mfe_response.mediaitemfaceembeddings()) {
      auto embedding = mfe.embedding();
      std::vector<float> mf_embedding;
      for (auto embed : embedding.embedding()) {
        mf_embedding.push_back(embed);
      }
      face_mediaitem_people_id.push_back(
          make_tuple(mfe.id(), mfe.mediaitemid(), mfe.peopleid()));
      embeddings.push_back(mf_embedding);
    }
    SPDLOG_DEBUG("user: {} faces: {}", user_id, faces.size());
    SPDLOG_INFO("running clustering algorithm");
  }
}

std::shared_ptr<Faces> Init(const std::string& models_dir,
                            const ComponentConfig& config,
                            std::shared_ptr<APIClient> api_client) {
  if (config.source == "onnx") {
    simdjson::ondemand::parser parser;
    simdjson::padded_string padded_config =
        simdjson::padded_string(config.params);
    simdjson::ondemand::document doc = parser.iterate(padded_config);
    float detection_threshold = 0.8f;
    if (doc["detection_threshold"].error() == simdjson::SUCCESS) {
      detection_threshold =
          static_cast<float>(doc["detection_threshold"].get_double());
    }
    std::string detection_model = "faces_det/scrfd_2.5g.onnx";
    if (doc["detection_threshold"].error() == simdjson::SUCCESS) {
      detection_model =
          std::string(doc["detection_model"].get_string().value());
    }
    std::string recognition_model = "faces_rec/webface_r50.onnx";
    if (doc["recognition_model"].error() == simdjson::SUCCESS) {
      recognition_model =
          std::string(doc["recognition_model"].get_string().value());
    }
    std::string clustering_lib = "faiss";
    if (doc["clustering_lib"].error() == simdjson::SUCCESS) {
      clustering_lib = std::string(doc["clustering_lib"].get_string().value());
    }
    return std::make_shared<ONNX>(
        std::make_shared<ONNXModel>(models_dir + "/" + detection_model,
                                    detection_threshold,
                                    models_dir + "/" + recognition_model),
        clustering_lib == "faiss" ? std::make_shared<FAISSClustering>()
                                  : nullptr,
        api_client);
  }
  return nullptr;
}

} // namespace faces

} // namespace components
