// Copyright 2025 Omkar Prabhu
#include "worker/faces.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"

using components::ComponentConfig;
using components::faces::Faces;
using components::faces::ModelInferenceInterface;
using components::faces::ONNX;
using components::faces::ONNXModel;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

class MockModelInference : public ModelInferenceInterface {
 public:
  MOCK_METHOD((std::vector<std::pair<std::string, std::vector<float>>>), Run,
              (const std::string& file_path), (override));
};

void assertFacesResult(std::unordered_map<std::string, std::string> expected,
                       std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(FacesTest, Init) {
  spdlog::set_level(spdlog::level::off);
  auto faces =
      components::faces::Init(ComponentConfig("onnx", "params"), nullptr);
  ASSERT_TRUE(faces != nullptr);
  auto onnx = std::dynamic_pointer_cast<ONNX>(faces);
  ASSERT_TRUE(onnx != nullptr);
  faces =
      components::faces::Init(ComponentConfig("unknown", "params"), nullptr);
  ASSERT_TRUE(faces == nullptr);
}

TEST(FacesTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  Faces faces;
  std::unordered_map<std::string, std::string> result =
      faces.Extract("", "", "", "");
  assertFacesResult({}, result);
}

TEST(FacesONNXTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockModelInference> mock_model =
      std::make_shared<MockModelInference>();
  ONNX onnx(mock_model, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      onnx.Extract("", "", "", "");
  assertFacesResult({}, result);
}

TEST(FacesONNXTest, Error) {
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockModelInference> mock_model =
      std::make_shared<MockModelInference>();
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillOnce(::testing::Throw(std::runtime_error("some error")));
  ONNX onnx(mock_model, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      onnx.Extract("", "", "", "file_path");
  assertFacesResult({}, result);
}

TEST(FacesONNXTest, Success) {
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockModelInference> mock_model =
      std::make_shared<MockModelInference>();
  std::vector<std::pair<std::string, std::vector<float>>> mock_response = {
      {"path/face1", {1.23, 4.56, 7.89}}, {"path/face2", {9.78, 6.54, 3.21}}};
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillOnce(::testing::Return(mock_response));
  ONNX onnx(mock_model, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      onnx.Extract("", "", "", "file_path");
  assertFacesResult(
      {
          {"faces", "path/face1,path/face2"},
      },
      result);
}
