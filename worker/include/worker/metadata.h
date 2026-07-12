// Copyright 2025 Omkar Prabhu
#pragma once

#include <spdlog/spdlog.h>
#include <unistd.h>

#include <cstdio>
#include <memory>
#include <mutex>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace metadata {

class ExifToolClientInterface {
 public:
  virtual ~ExifToolClientInterface() = default;
  virtual std::unordered_map<std::string, std::string> Extract(
      const std::string& file_path) = 0;
};

class ExifToolClient : public ExifToolClientInterface {
 public:
  ExifToolClient();
  ~ExifToolClient() override;
  std::unordered_map<std::string, std::string> Extract(
      const std::string& file_path) override;

 private:
  int in_fd[2], out_fd[2];
  FILE* in_stream{nullptr};
  FILE* out_stream{nullptr};
  pid_t pid;
  std::mutex io_mutex;
  std::string exiftool_path = "exiftool";
};

class Metadata {
 public:
  explicit Metadata(std::shared_ptr<ExifToolClientInterface> exif_tool_client,
                    std::shared_ptr<APIClient> api_client)
      : exif_tool_client_(exif_tool_client), api_client_(api_client) {}
  std::unordered_map<std::string, std::string> Extract(
      const std::string& id, const std::string& user_id,
      const std::string& mediaitem_id, const std::string& file_path);

 private:
  std::shared_ptr<ExifToolClientInterface> exif_tool_client_;
  std::shared_ptr<APIClient> api_client_;
};

std::string GetValue(const std::unordered_map<std::string, std::string>& data,
                     const std::vector<std::string>& keys);

std::string GetCoordinates(const std::string& location);

std::shared_ptr<Metadata> Init(const ComponentConfig& config,
                               std::shared_ptr<APIClient> api_client);
} // namespace metadata

} // namespace components
