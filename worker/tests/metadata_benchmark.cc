// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "metadata_test.cc" // NOLINT
#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"
#include "worker/metadata.h"

using components::ComponentConfig;
using components::metadata::ExifToolClientInterface;
using components::metadata::Metadata;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

static void BM_MetadataEmptyData(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockExifToolClient> mock_exif_client =
      std::make_shared<MockExifToolClient>();
  std::unordered_map<std::string, std::string> mock_data;
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  Metadata metadata(mock_exif_client, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        metadata.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_MetadataError(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockExifToolClient> mock_exif_client =
      std::make_shared<MockExifToolClient>();
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillRepeatedly(::testing::Throw(std::runtime_error("some error")));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  Metadata metadata(mock_exif_client, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        metadata.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_MetadataSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockExifToolClient> mock_exif_client =
      std::make_shared<MockExifToolClient>();
  std::unordered_map<std::string, std::string> mock_data = {
      {"VideoFrameRate", "30"},
      {"GPSLatitude", "19 deg 13' 11.99\" N"},
      {"GPSLongitude", "73 deg 6' 19.19\" E"},
      {"ExifImageWidth", "4032"},
      {"ImageWidth", "226"},
      {"ExifImageHeight", "3024"},
      {"ImageHeight", "4032"},
      {"Make", "Apple"},
      {"Model", "iPhone 15 Pro"},
      {"FocalLength", "4.2 mm"},
      {"FNumber", "1.6"},
      {"ISO", "640"},
      {"Megapixels", "24"},
      {"MIMEType", "image/heic"},
      {"ExposureTime", "1/25"},
      {"LivePhotoVideoIndex", "1112547328"},
      {"DateCreated", "2022:04:03 12:56:11"}};
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  Metadata metadata(mock_exif_client, mock_api_client);
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        metadata.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_MetadataEmptyData)->ThreadPerCpu();
BENCHMARK(BM_MetadataError)->ThreadPerCpu();
BENCHMARK(BM_MetadataSuccess)->ThreadPerCpu();
