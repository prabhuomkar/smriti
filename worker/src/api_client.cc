// Copyright 2025 Omkar Prabhu
#include "worker/api_client.h"

#include <google/protobuf/empty.pb.h>
#include <grpcpp/grpcpp.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>

#include "protos/api.grpc.pb.h"
#include "protos/api.pb.h"

namespace services {

namespace api {

APIClient::APIClient(std::shared_ptr<grpc::Channel> channel)
    : stub_(API::NewStub(channel)) {}

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

} // namespace api

} // namespace services
