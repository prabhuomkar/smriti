// Copyright 2025 Omkar Prabhu
#include "worker/api_client.h"

#include <gmock/gmock.h>
#include <grpcpp/grpcpp.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <utility>

#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"

using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

class APIClientTest : public ::testing::Test {
 protected:
  void SetUp() override {
    spdlog::set_level(spdlog::level::off);
    auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
    mock_stub_ = mock_api_stub.get();
    client_ =
        std::make_unique<APIClient>(std::move(mock_api_stub)); // move ownership
  }

  MockAPIStub* mock_stub_;
  std::unique_ptr<APIClient> client_;
};

TEST_F(APIClientTest, GetWorkerConfigSuccess) {
  EXPECT_CALL(*mock_stub_, GetWorkerConfig(_, _, _))
      .WillOnce(Invoke([&](grpc::ClientContext*, const google::protobuf::Empty&,
                           ConfigResponse* resp) {
        ConfigResponse response;
        response.set_config("worker-config");
        *resp = response;
        return grpc::Status::OK;
      }));
  std::string config = client_->GetWorkerConfig();
  EXPECT_EQ(config, "worker-config");
}

TEST_F(APIClientTest, GetWorkerConfigFailure) {
  EXPECT_CALL(*mock_stub_, GetWorkerConfig(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  std::string config = client_->GetWorkerConfig();
  EXPECT_EQ(config, "");
}

TEST_F(APIClientTest, GetMediaItemProcessSuccess) {
  EXPECT_CALL(*mock_stub_, GetMediaItemProcess(_, _, _))
      .WillOnce(Invoke([&](grpc::ClientContext*, const google::protobuf::Empty&,
                           MediaItemProcessResponse* resp) {
        MediaItemProcessResponse response;
        response.set_id("id");
        response.set_mediaitemid("mediaitem-id");
        response.add_components(MediaItemComponent::METADATA);
        response.add_components(MediaItemComponent::PREVIEW_THUMBNAIL);
        (*response.mutable_payload())["key"] = "value";
        *resp = response;
        return grpc::Status::OK;
      }));
  MediaItemProcessResponse mediaitem_process = client_->GetMediaItemProcess();
  EXPECT_EQ(mediaitem_process.id(), "id");
  EXPECT_EQ(mediaitem_process.mediaitemid(), "mediaitem-id");
  ASSERT_EQ(mediaitem_process.components_size(), 2);
  EXPECT_EQ(mediaitem_process.components(0), MediaItemComponent::METADATA);
  EXPECT_EQ(mediaitem_process.components(1),
            MediaItemComponent::PREVIEW_THUMBNAIL);
  ASSERT_EQ(mediaitem_process.payload().size(), 1);
  EXPECT_EQ(mediaitem_process.payload().at("key"), "value");
}

TEST_F(APIClientTest, GetMediaItemProcessFailure) {
  EXPECT_CALL(*mock_stub_, GetMediaItemProcess(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemProcessResponse mediaitem_process = client_->GetMediaItemProcess();
  EXPECT_TRUE(mediaitem_process.id().empty());
  EXPECT_TRUE(mediaitem_process.mediaitemid().empty());
  EXPECT_EQ(mediaitem_process.components_size(), 0);
  EXPECT_EQ(mediaitem_process.payload().size(), 0);
}

TEST_F(APIClientTest, SaveMediaItemMetadataSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemMetadata(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemMetadataRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemMetadataRequest request;
  bool ok = client_->SaveMediaItemMetadata(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemMetadataFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemMetadata(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemMetadataRequest request;
  bool ok = client_->SaveMediaItemMetadata(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPreviewThumbnailSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillOnce(Invoke(
          [&](grpc::ClientContext*, const MediaItemPreviewThumbnailRequest&,
              google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemPreviewThumbnailRequest request;
  bool ok = client_->SaveMediaItemPreviewThumbnail(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPreviewThumbnailFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPreviewThumbnail(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemPreviewThumbnailRequest request;
  bool ok = client_->SaveMediaItemPreviewThumbnail(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPlaceSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPlace(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemPlaceRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemPlaceRequest request;
  bool ok = client_->SaveMediaItemPlace(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPlaceFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPlace(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemPlaceRequest request;
  bool ok = client_->SaveMediaItemPlace(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemThingSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemThing(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemThingRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemThingRequest request;
  bool ok = client_->SaveMediaItemThing(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemThingFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemThing(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemThingRequest request;
  bool ok = client_->SaveMediaItemThing(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemFacesSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemFaces(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemFacesRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemFacesRequest request;
  bool ok = client_->SaveMediaItemFaces(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemFacesFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemFaces(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemFacesRequest request;
  bool ok = client_->SaveMediaItemFaces(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPeopleSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPeople(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemPeopleRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemPeopleRequest request;
  bool ok = client_->SaveMediaItemPeople(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemPeopleFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemPeople(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemPeopleRequest request;
  bool ok = client_->SaveMediaItemPeople(request);
  EXPECT_FALSE(ok);
}

TEST_F(APIClientTest, SaveMediaItemFinalResultSuccess) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemFinalResult(_, _, _))
      .WillOnce(
          Invoke([&](grpc::ClientContext*, const MediaItemFinalResultRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  MediaItemFinalResultRequest request;
  bool ok = client_->SaveMediaItemFinalResult(request);
  EXPECT_TRUE(ok);
}

TEST_F(APIClientTest, SaveMediaItemFinalResultFailure) {
  EXPECT_CALL(*mock_stub_, SaveMediaItemFinalResult(_, _, _))
      .WillOnce(Return(
          grpc::Status(grpc::StatusCode::UNAVAILABLE, "Service unavailable")));
  MediaItemFinalResultRequest request;
  bool ok = client_->SaveMediaItemFinalResult(request);
  EXPECT_FALSE(ok);
}
