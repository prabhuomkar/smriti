// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
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

#include "faces_test.cc" // NOLINT
#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"
#include "worker/faces.h"

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

static void BM_FacesEmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  Faces faces;
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        faces.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_FacesONNXEmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockFacesModelInference> mock_model =
      std::make_shared<MockFacesModelInference>();
  ONNX onnx(mock_model, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        onnx.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_FacesONNXError(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockFacesModelInference> mock_model =
      std::make_shared<MockFacesModelInference>();
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillRepeatedly(::testing::Throw(std::runtime_error("some error")));
  ONNX onnx(mock_model, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        onnx.Extract("", "", "", "file_path");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_FacesONNXSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemFaces(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockFacesModelInference> mock_model =
      std::make_shared<MockFacesModelInference>();
  std::vector<std::pair<std::string, std::vector<float>>> mock_response = {
      {"path/face1", {1.23, 4.56, 7.89}}, {"path/face2", {9.78, 6.54, 3.21}}};
  EXPECT_CALL(*mock_model, Run(::testing::_))
      .WillRepeatedly(::testing::Return(mock_response));
  ONNX onnx(mock_model, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        onnx.Extract("", "", "", "file_path");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_FacesEmptyInput)->ThreadPerCpu();
BENCHMARK(BM_FacesONNXEmptyInput)->ThreadPerCpu();
BENCHMARK(BM_FacesONNXError)->ThreadPerCpu();
BENCHMARK(BM_FacesONNXSuccess)->ThreadPerCpu();
