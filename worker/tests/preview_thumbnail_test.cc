// Copyright 2025 Omkar Prabhu
#include "worker/preview_thumbnail.h"

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
using components::previewthumbnail::ImageConverterClientInterface;
using components::previewthumbnail::PreviewThumbnail;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

class MockImageConverterClient : public ImageConverterClientInterface {
 public:
  MOCK_METHOD((std::string), Convert,
              (const std::string& input_file_path,
               const std::string& output_file_path, int image_size),
              (override));
};

void assertPreviewThumbnailResult(
    std::unordered_map<std::string, std::string> expected,
    std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(PreviewThumbnailTest, Init) {
  spdlog::set_level(spdlog::level::off);
  auto previewthumbnail = components::previewthumbnail::Init(
      ComponentConfig("",
                      "{\"image_quality\":50,\"thumbnail_size\":256,"
                      "\"placeholder_size\":2}"),
      nullptr);
  ASSERT_TRUE(previewthumbnail != nullptr);
}

TEST(PreviewThumbnailTest, Error) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockImageConverterClient> mock_image_converter_client =
      std::make_shared<MockImageConverterClient>();
  EXPECT_CALL(*mock_image_converter_client,
              Convert(::testing::_, ::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Throw(std::runtime_error("some error")));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillOnce(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  PreviewThumbnail previewthumbnail(mock_image_converter_client,
                                    mock_api_client);
  std::unordered_map<std::string, std::string> result =
      previewthumbnail.Generate("", "", "", "",
                                MediaItemType_Name(MediaItemType::PHOTO));
  assertPreviewThumbnailResult({{"status", "FAILED"},
                                {"preview_url", ""},
                                {"thumbnail_url", ""},
                                {"placeholder", ""}},
                               result);
}

TEST(PreviewThumbnailTest, PhotoSuccess) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockImageConverterClient> mock_image_converter_client =
      std::make_shared<MockImageConverterClient>();
  std::string mock_data = "path-kind";
  EXPECT_CALL(*mock_image_converter_client,
              Convert(::testing::_, ::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillOnce(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  PreviewThumbnail previewthumbnail(mock_image_converter_client,
                                    mock_api_client);
  std::unordered_map<std::string, std::string> result =
      previewthumbnail.Generate("", "", "", "",
                                MediaItemType_Name(MediaItemType::PHOTO));
  assertPreviewThumbnailResult({{"preview_url", "path-kind"},
                                {"thumbnail_url", "path-kind"},
                                {"placeholder", "path-kind"},
                                {"status", "READY"}},
                               result);
}

TEST(PreviewThumbnailTest, VideoSuccess) {
  spdlog::set_level(spdlog::level::off);
  std::shared_ptr<MockImageConverterClient> mock_image_converter_client =
      std::make_shared<MockImageConverterClient>();
  std::string mock_data = "path-kind";
  EXPECT_CALL(*mock_image_converter_client,
              Convert(::testing::_, ::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Return(mock_data));
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillOnce(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  PreviewThumbnail previewthumbnail(mock_image_converter_client,
                                    mock_api_client);
  std::unordered_map<std::string, std::string> result =
      previewthumbnail.Generate("", "", "", "",
                                MediaItemType_Name(MediaItemType::VIDEO));
  assertPreviewThumbnailResult({{"preview_url", ""},
                                {"thumbnail_url", ""},
                                {"placeholder", ""},
                                {"status", "READY"}},
                               result);
}
