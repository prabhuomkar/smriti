// Copyright 2025 Omkar Prabhu
#include <simdjson.h>
#include <spdlog/sinks/stdout_sinks.h>
#include <spdlog/spdlog.h>

#include <atomic>
#include <csignal>
#include <cstdlib>
#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>
#include <vector>

#ifndef DEFAULT_VERSION
#define DEFAULT_VERSION "dev"
#endif
#ifndef DEFAULT_GIT_SHA
#define DEFAULT_GIT_SHA "-"
#endif

#include "protos/api.pb.h"
#include "worker/api_client.h"
#include "worker/components.h"
#include "worker/config.h"
#include "worker/metadata.h"
#include "worker/places.h"

using components::ComponentConfig;
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
  std::unordered_map<std::string, ComponentConfig> component_configs =
      components::ParseComponentConfig(worker_config);
  std::shared_ptr<components::metadata::Metadata> metadata_component;
  std::shared_ptr<components::places::Places> places_component;
  for (const auto& [name, config] : component_configs) {
    if (name == MediaItemComponent_Name(MediaItemComponent::METADATA)) {
      metadata_component = components::metadata::Init(config);
    } else if (name == MediaItemComponent_Name(MediaItemComponent::PLACES)) {
      places_component = components::places::Init(config);
    } else if (name ==
               MediaItemComponent_Name(MediaItemComponent::PREVIEW_THUMBNAIL)) {
      // TODO(omkar): initialize this component
    } else if (name ==
               MediaItemComponent_Name(MediaItemComponent::CLASSIFICATION)) {
      // TODO(omkar): initialize this component

    } else if (name == MediaItemComponent_Name(MediaItemComponent::FACES)) {
      // TODO(omkar): initialize this component

    } else if (name == MediaItemComponent_Name(MediaItemComponent::OCR)) {
      // TODO(omkar): initialize this component

    } else if (name == MediaItemComponent_Name(MediaItemComponent::SEARCH)) {
      // TODO(omkar): initialize this component
    }
  }

  while (!terminating) {
    spdlog::info("worker running");
    sleep(10);
    // TODO(omkar): fetch mediaitem to process
  }

  return 0;
}
