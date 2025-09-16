// Copyright 2025 Omkar Prabhu
#include "worker/metadata.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"

using components::ComponentConfig;
using components::metadata::ExifToolClientInterface;
using components::metadata::Metadata;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

class MockExifToolClient : public ExifToolClientInterface {
 public:
  MOCK_METHOD((std::unordered_map<std::string, std::string>), Extract,
              (const std::string& file_path), (override));
};

void assertMetadataResult(std::unordered_map<std::string, std::string> expected,
                          std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(MetadataTest, Init) {
  spdlog::set_level(spdlog::level::off);
  auto metadata = components::metadata::Init(ComponentConfig("", ""), nullptr);
  ASSERT_TRUE(metadata != nullptr);
}

TEST(MetadataTest, EmptyData) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockExifToolClient> mock_exif_client =
      std::make_shared<MockExifToolClient>();
  std::unordered_map<std::string, std::string> mock_data;
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillOnce(::testing::Return(mock_data));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  Metadata metadata(mock_exif_client, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      metadata.Extract("", "", "", "");
  assertMetadataResult(
      {
          {"status", "PROCESSING"},
          {"type", "unknown"},
          {"category", "default"},
          {"latitude", ""},
          {"longitude", ""},
          {"fps", ""},
          {"height", ""},
          {"width", ""},
          {"camera_make", ""},
          {"camera_model", ""},
          {"focal_length", ""},
          {"aperture_fnumber", ""},
          {"iso_equivalent", ""},
          {"exposure_time", ""},
          {"mime_type", ""},
          {"creation_time", ""},
          {"megapixels", ""},
          {"exifdata", "{}"},
      },
      result);
}

TEST(MetadataTest, Failure) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockExifToolClient> mock_exif_client =
      std::make_shared<MockExifToolClient>();
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillOnce(::testing::Throw(std::runtime_error("some error")));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  Metadata metadata(mock_exif_client, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      metadata.Extract("", "", "", "");
  assertMetadataResult(
      {{"status", "FAILED"}, {"type", "unknown"}, {"category", "default"}},
      result);
}

TEST(MetadataTest, Success) {
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
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemMetadata(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  EXPECT_CALL(*mock_exif_client, Extract(::testing::_))
      .WillOnce(::testing::Return(mock_data));
  Metadata metadata(mock_exif_client, mock_api_client);
  std::unordered_map<std::string, std::string> result =
      metadata.Extract("", "", "", "");
  assertMetadataResult(
      {{"status", "PROCESSING"},
       {"type", "photo"},
       {"category", "live"},
       {"latitude", "19.219997"},
       {"longitude", "73.105331"},
       {"fps", "30"},
       {"height", "3024"},
       {"width", "4032"},
       {"camera_make", "Apple"},
       {"camera_model", "iPhone 15 Pro"},
       {"focal_length", "4.2 mm"},
       {"aperture_fnumber", "1.6"},
       {"iso_equivalent", "640"},
       {"exposure_time", "1/25"},
       {"mime_type", "image/heic"},
       {"megapixels", "24"},
       {"creation_time", "2022-04-03 12:56:11"},
       {"exifdata",
        "{\"VideoFrameRate\":\"30\",\"ExifImageWidth\":"
        "\"4032\",\"GPSLatitude\":\"19 deg 13' 11.99\" "
        "N\",\"GPSLongitude\":\"73 deg 6' 19.19\" "
        "E\",\"ImageWidth\":\"226\",\"Make\":\"Apple\","
        "\"Model\":\"iPhone 15 Pro\",\"FocalLength\":\"4.2 "
        "mm\",\"ExifImageHeight\":\"3024\",\"FNumber\":\"1."
        "6\",\"ImageHeight\":\"4032\",\"ISO\":\"640\",\"Megapixels\":\"24\","
        "\"MIMEType\":\"image/heic\",\"ExposureTime\":\"1/"
        "25\",\"LivePhotoVideoIndex\":\"1112547328\","
        "\"DateCreated\":\"2022:04:03 12:56:11\"}"}},
      result);
}
