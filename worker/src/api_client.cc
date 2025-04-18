// Copyright 2025 Omkar Prabhu
#include "worker/api_client.h"

#include <google/protobuf/empty.pb.h>
#include <grpcpp/grpcpp.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <utility>

#include "protos/api.grpc.pb.h"
#include "protos/api.pb.h"

namespace services {

namespace api {

APIClient::APIClient(std::shared_ptr<grpc::Channel> channel)
    : stub_(API::NewStub(channel)) {}

APIClient::APIClient(std::unique_ptr<API::StubInterface> stub)
    : stub_(std::move(stub)) {}

std::string APIClient::GetWorkerConfig() {
  google::protobuf::Empty request;
  ConfigResponse response;
  grpc::ClientContext context;
  grpc::Status status = stub_->GetWorkerConfig(&context, request, &response);
  if (!status.ok()) {
    spdlog::error("error getting worker config: {}", status.error_message());
    return "";
  }
  return response.config();
}

MediaItemProcessResponse APIClient::GetMediaItemProcess() {
  google::protobuf::Empty request;
  MediaItemProcessResponse response;
  grpc::ClientContext context;
  grpc::Status status =
      stub_->GetMediaItemProcess(&context, request, &response);
  if (!status.ok()) {
    spdlog::error("error getting media item process: {}",
                  status.error_message());
    return {};
  }
  return response;
}

bool APIClient::SaveMediaItemMetadata(const MediaItemMetadataRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status =
      stub_->SaveMediaItemMetadata(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemPreviewThumbnail(
    const MediaItemPreviewThumbnailRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status =
      stub_->SaveMediaItemPreviewThumbnail(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemPlace(const MediaItemPlaceRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status = stub_->SaveMediaItemPlace(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemThing(const MediaItemThingRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status = stub_->SaveMediaItemThing(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemFaces(const MediaItemFacesRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status = stub_->SaveMediaItemFaces(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemPeople(const MediaItemPeopleRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status =
      stub_->SaveMediaItemPeople(&context, request, &response);
  return status.ok();
}

bool APIClient::SaveMediaItemFinalResult(
    const MediaItemFinalResultRequest& request) {
  google::protobuf::Empty response;
  grpc::ClientContext context;
  grpc::Status status =
      stub_->SaveMediaItemFinalResult(&context, request, &response);
  return status.ok();
}

} // namespace api

} // namespace services
