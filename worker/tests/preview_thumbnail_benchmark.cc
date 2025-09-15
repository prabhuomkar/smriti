// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "preview_thumbnail_test.cc" // NOLINT
#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"
#include "worker/preview_thumbnail.h"

using components::ComponentConfig;
using components::previewthumbnail::ImageConverterClientInterface;
using components::previewthumbnail::PreviewThumbnail;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

static void BM_PreviewThumbnailInit(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  ComponentConfig config =
      ComponentConfig("",
                      "{\"image_quality\":50,\"thumbnail_size\":256,"
                      "\"placeholder_size\":2}");
  for (auto _ : state) {
    auto previewthumbnail = components::previewthumbnail::Init(
        ComponentConfig("",
                        "{\"image_quality\":50,\"thumbnail_size\":256,"
                        "\"placeholder_size\":2}"),
        nullptr);
    benchmark::DoNotOptimize(previewthumbnail);
  }
}

static void BM_PreviewThumbnailPhotoSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillRepeatedly(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockImageConverterClient> mock_image_converter_client =
      std::make_shared<MockImageConverterClient>();
  std::string mock_data = "path-kind";
  EXPECT_CALL(*mock_image_converter_client,
              Convert(::testing::_, ::testing::_, ::testing::_, ::testing::_,
                      ::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  for (auto _ : state) {
    PreviewThumbnail previewthumbnail(mock_image_converter_client,
                                      mock_api_client);
    std::unordered_map<std::string, std::string> output =
        previewthumbnail.Generate("", "", "", "", "photo");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_PreviewThumbnailVideoSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillRepeatedly(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockImageConverterClient> mock_image_converter_client =
      std::make_shared<MockImageConverterClient>();
  std::string mock_data = "path-kind";
  EXPECT_CALL(*mock_image_converter_client,
              Convert(::testing::_, ::testing::_, ::testing::_, ::testing::_,
                      ::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  for (auto _ : state) {
    PreviewThumbnail previewthumbnail(mock_image_converter_client,
                                      mock_api_client);
    std::unordered_map<std::string, std::string> output =
        previewthumbnail.Generate("", "", "", "", "video");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_PreviewThumbnailInit)->ThreadPerCpu();
BENCHMARK(BM_PreviewThumbnailPhotoSuccess)->ThreadPerCpu();
BENCHMARK(BM_PreviewThumbnailVideoSuccess)->ThreadPerCpu();
