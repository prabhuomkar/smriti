// Copyright 2025 Omkar Prabhu
#pragma once

#include <grpcpp/grpcpp.h>

#include <memory>
#include <string>

#include "protos/api.grpc.pb.h"
#include "protos/api.pb.h"

namespace services {

namespace api {

class APIClient {
 public:
  explicit APIClient(std::shared_ptr<grpc::Channel> channel);
  explicit APIClient(std::unique_ptr<API::StubInterface> stub);
  std::string GetWorkerConfig();
  MediaItemProcessResponse GetMediaItemProcess();
  UsersResponse GetUsers();
  MediaItemFaceEmbeddingsResponse GetMediaItemFaceEmbeddings(
      const MediaItemFaceEmbeddingsRequest& request);
  bool SaveMediaItemMetadata(const MediaItemMetadataRequest& request);
  bool SaveMediaItemPreviewThumbnail(
      const MediaItemPreviewThumbnailRequest& request);
  bool SaveMediaItemPlace(const MediaItemPlaceRequest& request);
  bool SaveMediaItemFaces(const MediaItemFacesRequest& request);
  bool SaveMediaItemPeople(const MediaItemPeopleRequest& request);
  bool SaveMediaItemFinalResult(const MediaItemFinalResultRequest& request);

 private:
  std::unique_ptr<API::StubInterface> stub_;
};

} // namespace api

} // namespace services
