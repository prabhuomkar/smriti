// Copyright 2025 Omkar Prabhu
#include <spdlog/sinks/stdout_sinks.h>
#include <spdlog/spdlog.h>

#include <atomic>
#include <csignal>
#include <cstdlib>
#include <iostream>
#include <memory>
#include <string>

#ifndef DEFAULT_VERSION
#define DEFAULT_VERSION "dev"
#endif
#ifndef DEFAULT_GIT_SHA
#define DEFAULT_GIT_SHA "-"
#endif

#include "worker/config.h"

constexpr const char kVersion[] = DEFAULT_VERSION;
constexpr const char kGitSha[] = DEFAULT_GIT_SHA;

std::atomic<bool> terminating(false);

void gracefulShutdown(int signum) {
  spdlog::info("stopping worker");
  // TODO(omkar): clean up gRPC client
  terminating = true;
}

int main() {
  std::cout << "Version: " << kVersion << std::endl;
  std::cout << "Git SHA: " << kGitSha << std::endl;

  auto cfg = std::make_unique<Config>();

  auto stdout_sink = std::make_shared<spdlog::sinks::stdout_sink_mt>();
  auto logger = std::make_shared<spdlog::logger>("", stdout_sink);
  logger->set_level(cfg->log_level);
  logger->set_pattern("%Y/%m/%d %H:%M:%S %l %v");
  spdlog::set_default_logger(logger);

  std::signal(SIGTERM, gracefulShutdown);
  std::signal(SIGINT, gracefulShutdown);

  while (!terminating) {
    // TODO(omkar): Initialize gRPC client
  }

  return 0;
}
