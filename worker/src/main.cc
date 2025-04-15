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

#include "worker/api_client.h"
#include "worker/config.h"

using services::api::APIClient;

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

  APIClient api_client(grpc::CreateChannel(cfg->api_host + ":" + cfg->api_port,
                                           grpc::InsecureChannelCredentials()));
  std::string worker_config = api_client.GetWorkerConfig();
  spdlog::info("worker config: {}", worker_config);

  while (!terminating) {
    // TODO(omkar): Pull jobs from API server and execute graph
    spdlog::info("worker running");
    std::this_thread::sleep_for(std::chrono::seconds(10));
  }

  return 0;
}
