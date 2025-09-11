// Copyright 2025 Omkar Prabhu
#pragma once

#include <unistd.h>

#include <cstdio>
#include <memory>
#include <mutex>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/components.h"

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
  ExifToolClient() {
    if (pipe(in_fd) == -1 || pipe(out_fd) == -1) {
      throw std::runtime_error("pipe failed");
    }

    pid = fork();
    if (pid == -1) {
      throw std::runtime_error("fork failed");
    }

    if (pid == 0) {
      dup2(in_fd[0], STDIN_FILENO);
      dup2(out_fd[1], STDOUT_FILENO);
      close(in_fd[1]);
      close(out_fd[0]);
      execlp(exiftool_path.c_str(), exiftool_path.c_str(), "-stay_open", "True",
             "-@", "-", static_cast<char*>(nullptr));
      _exit(1);
    }

    close(in_fd[0]);
    close(out_fd[1]);

    in_stream = fdopen(in_fd[1], "w");
    out_stream = fdopen(out_fd[0], "r");
    if (!in_stream || !out_stream) {
      throw std::runtime_error("fdopen failed");
    }
  }

  ~ExifToolClient() {
    if (in_stream) {
      fprintf(in_stream, "-stay_open\nFalse\n");
      fflush(in_stream);
      fclose(in_stream);
    }
    if (out_stream) {
      fclose(out_stream);
    }
    waitpid(pid, nullptr, 0);
  }

  std::unordered_map<std::string, std::string> Extract(
      const std::string& file_path) override {
    std::lock_guard<std::mutex> lock(io_mutex);

    fprintf(in_stream, "-s\n-s\n%s\n-execute\n", file_path.c_str());
    fflush(in_stream);

    std::unordered_map<std::string, std::string> result;
    char buffer[4096];
    while (fgets(buffer, sizeof(buffer), out_stream)) {
      std::string line(buffer);
      if (line.find("{ready}") != std::string::npos) {
        break;
      }
      auto pos = line.find_first_of(':');
      if (pos != std::string::npos) {
        std::string key = line.substr(0, pos);
        std::string value = line.substr(pos + 2, line.length() - pos - 3);
        result[key] = value;
      }
    }

    return result;
  }

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
  explicit Metadata(std::shared_ptr<ExifToolClientInterface> exif_tool_client)
      : exif_tool_client_(exif_tool_client) {}
  std::unordered_map<std::string, std::string> Extract(
      const std::string& user_id, const std::string& mediaitem_id,
      const std::string& file_path);

 private:
  std::shared_ptr<ExifToolClientInterface> exif_tool_client_;
};

std::shared_ptr<Metadata> Init(const ComponentConfig& config);
} // namespace metadata

} // namespace components
