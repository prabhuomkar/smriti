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
  std::string GetWorkerConfig();

 private:
  std::unique_ptr<API::Stub> stub_;
};

} // namespace api

} // namespace services
