// Copyright 2025 Omkar Prabhu
#pragma once

#include <spdlog/spdlog.h>

#include <cstdlib>
#include <iostream>
#include <memory>
#include <string>

class Config {
 public:
  Config() {
    const char* env_log_level = std::getenv("SMRITI_WORKER_LOG_LEVEL");
    log_level = spdlog::level::info;
    if (env_log_level) {
      std::string level_str(env_log_level);
      if (level_str == "debug")
        log_level = spdlog::level::debug;
      else if (level_str == "info")
        log_level = spdlog::level::info;
      else if (level_str == "warn")
        log_level = spdlog::level::warn;
      else if (level_str == "error")
        log_level = spdlog::level::err;
      else if (level_str == "critical")
        log_level = spdlog::level::critical;
    }

    const char* env_port = std::getenv("SMRITI_WORKER_PORT");
    port = (env_port && std::strlen(env_port) > 0) ? std::string(env_port)
                                                   : "15002";

    const char* env_api_host = std::getenv("SMRITI_API_HOST");
    api_host = (env_api_host && std::strlen(env_api_host) > 0)
                   ? std::string(env_api_host)
                   : "127.0.0.1";

    const char* env_api_port = std::getenv("SMRITI_API_PORT");
    api_port = (env_api_port && std::strlen(env_api_port) > 0)
                   ? std::string(env_api_port)
                   : "15001";

    const char* env_models_dir = std::getenv("SMRITI_MODELS_DIR");
    models_dir = (env_models_dir && std::strlen(env_models_dir) > 0)
                     ? std::string(env_models_dir)
                     : "";
  }

  spdlog::level::level_enum log_level;
  std::string port;
  std::string api_host;
  std::string api_port;
  std::string models_dir;
};
