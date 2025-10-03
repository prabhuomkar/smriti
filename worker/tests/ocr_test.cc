// Copyright 2025 Omkar Prabhu
#include "worker/ocr.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"

using components::ComponentConfig;
using components::ocr::ModelInferenceInterface;
using components::ocr::OCR;
using components::ocr::PaddlePaddle;
using components::ocr::PaddlePaddleModel;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

class MockOCRModelInference : public ModelInferenceInterface {
 public:
  MOCK_METHOD((std::vector<std::pair<std::string, float>>), Run,
              (const std::string& file_path), (override));
};

void assertOCRResult(std::unordered_map<std::string, std::string> expected,
                     std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(OCRTest, Init) {
  spdlog::set_level(spdlog::level::off);
  auto ocr = components::ocr::Init(
      "../../../models", ComponentConfig("paddlepaddle", "params"), nullptr);
  ASSERT_TRUE(ocr != nullptr);
  auto paddlepaddle = std::dynamic_pointer_cast<PaddlePaddle>(ocr);
  ASSERT_TRUE(paddlepaddle != nullptr);
  ocr = components::ocr::Init("../../../models",
                              ComponentConfig("unknown", "params"), nullptr);
  ASSERT_TRUE(ocr == nullptr);
}

TEST(OCRTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  components::ocr::OCR ocr;
  std::unordered_map<std::string, std::string> result =
      ocr.Extract("", "", "", "");
  assertOCRResult({}, result);
}

TEST(OCRPaddlePaddleTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockOCRModelInference> mock_model =
      std::make_shared<MockOCRModelInference>();
  PaddlePaddle paddlepaddle(mock_model, nullptr);
  std::unordered_map<std::string, std::string> result =
      paddlepaddle.Extract("", "", "", "");
  assertOCRResult({}, result);
}

TEST(OCRPaddlePaddleTest, Error) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockOCRModelInference> mock_model =
      std::make_shared<MockOCRModelInference>();
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillOnce(::testing::Throw(std::runtime_error("some error")));
  PaddlePaddle paddlepaddle(mock_model, nullptr);
  std::unordered_map<std::string, std::string> result =
      paddlepaddle.Extract("", "", "", "file_path");
  assertOCRResult({}, result);
}

TEST(OCRPaddlePaddleTest, Success) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockOCRModelInference> mock_model =
      std::make_shared<MockOCRModelInference>();
  std::vector<std::pair<std::string, float>> mock_response = {
      {"thats what she said", 99.86},
      {"boy have you lost your mind", 98.12},
      {"spiderface", 75.43}};
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillOnce(::testing::Return(mock_response));
  PaddlePaddle paddlepaddle(mock_model, nullptr);
  std::unordered_map<std::string, std::string> result =
      paddlepaddle.Extract("", "", "", "file_path");
  assertOCRResult(
      {
          {"detected_text",
           "thats what she said boy have you lost your mind spiderface"},
      },
      result);
}
