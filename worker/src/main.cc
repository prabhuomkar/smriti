// Copyright 2025 Omkar Prabhu
#include <Magick++.h>
#include <simdjson.h>
#include <spdlog/sinks/stdout_sinks.h>
#include <spdlog/spdlog.h>

#include <atomic>
#include <chrono>
#include <csignal>
#include <cstdlib>
#include <future>
#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>
#include <utility>
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
#include "worker/faces.h"
#include "worker/metadata.h"
#include "worker/ocr.h"
#include "worker/places.h"
#include "worker/preview_thumbnail.h"

using components::ComponentConfig;
using services::api::APIClient;

constexpr const char kVersion[] = DEFAULT_VERSION;
constexpr const char kGitSha[] = DEFAULT_GIT_SHA;

std::atomic<bool> terminating(false);

void gracefulShutdown(int signum) {
  SPDLOG_INFO("stopping worker");
  terminating = true;
}

std::future<void> ProcessMediaItem(
    MediaItemProcessResponse response, std::shared_ptr<APIClient> api_client,
    std::shared_ptr<components::metadata::Metadata> metadata_component,
    std::shared_ptr<components::previewthumbnail::PreviewThumbnail>
        previewthumbnail_component,
    std::shared_ptr<components::places::Places> places_component,
    std::shared_ptr<components::faces::Faces> faces_component,
    std::shared_ptr<components::ocr::OCR> ocr_component) {
  return std::async(
      std::launch::async,
      [response, api_client, metadata_component, previewthumbnail_component,
       places_component, faces_component, ocr_component]() {
        try {
          SPDLOG_INFO("processing mediaitem: {}", response.id());

          std::unordered_map<std::string, std::string> result;
          for (const auto& [key, value] : response.payload()) {
            result[key] = value;
          }

          for (auto component : response.components()) {
            SPDLOG_INFO("component: {}", MediaItemComponent_Name(component));
            std::unordered_map<std::string, std::string> component_result;

            if (component == MediaItemComponent::METADATA) {
              component_result = metadata_component->Extract(
                  response.id(), response.userid(), response.mediaitemid(),
                  result["source_url"]);
            } else if (component == MediaItemComponent::PREVIEW_THUMBNAIL) {
              component_result = previewthumbnail_component->Generate(
                  response.id(), response.userid(), response.mediaitemid(),
                  result["source_url"], result["type"]);
            } else if (component == MediaItemComponent::PLACES) {
              component_result = places_component->ReverseGeocode(
                  response.id(), response.userid(), response.mediaitemid(),
                  result["latitude"], result["longitude"]);
            } else if (component == MediaItemComponent::FACES) {
              component_result = faces_component->Extract(
                  response.id(), response.userid(), response.mediaitemid(),
                  result["preview_url"]);
            } else if (component == MediaItemComponent::OCR) {
              component_result = ocr_component->Extract(
                  response.id(), response.userid(), response.mediaitemid(),
                  result["preview_url"]);
            }

            for (const auto& [key, value] : component_result) {
              result[key] = value;
            }
          }

          MediaItemFinalResultRequest final_request;
          final_request.set_id(response.id());
          final_request.set_userid(response.userid());
          final_request.set_mediaitemid(response.mediaitemid());
          final_request.set_detectedtext(result.find("detected_text") !=
                                                 result.end()
                                             ? result["detected_text"]
                                             : "");
          final_request.set_caption("");
          api_client->SaveMediaItemFinalResult(final_request);

          faces_component->Cluster();
        } catch (const std::exception& e) {
          SPDLOG_ERROR("error processing mediaitem {}: {}", response.id(),
                       e.what());
        }
      });
}

int main(int argc, char** argv) {
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

  std::shared_ptr<APIClient> api_client = std::make_shared<APIClient>(
      grpc::CreateChannel(cfg->api_host + ":" + cfg->api_port,
                          grpc::InsecureChannelCredentials()));
  std::string worker_config = api_client->GetWorkerConfig();
  SPDLOG_INFO("worker config: {}", worker_config);

  std::shared_ptr<components::metadata::Metadata> metadata_component;
  std::shared_ptr<components::places::Places> places_component;
  std::shared_ptr<components::previewthumbnail::PreviewThumbnail>
      previewthumbnail_component;
  std::shared_ptr<components::faces::Faces> faces_component;
  std::shared_ptr<components::ocr::OCR> ocr_component;

  std::unordered_map<std::string, ComponentConfig> component_configs =
      components::ParseComponentConfig(worker_config);
  for (const auto& [name, component_config] : component_configs) {
    if (name == MediaItemComponent_Name(MediaItemComponent::METADATA)) {
      metadata_component =
          components::metadata::Init(component_config, api_client);
    } else if (name == MediaItemComponent_Name(MediaItemComponent::PLACES)) {
      places_component = components::places::Init(component_config, api_client);
    } else if (name ==
               MediaItemComponent_Name(MediaItemComponent::PREVIEW_THUMBNAIL)) {
      Magick::InitializeMagick(*argv);
      previewthumbnail_component =
          components::previewthumbnail::Init(component_config, api_client);
    } else if (name == MediaItemComponent_Name(MediaItemComponent::FACES)) {
      faces_component = components::faces::Init(cfg->models_dir,
                                                component_config, api_client);
    } else if (name == MediaItemComponent_Name(MediaItemComponent::OCR)) {
      ocr_component =
          components::ocr::Init(cfg->models_dir, component_config, api_client);
    } else if (name == MediaItemComponent_Name(MediaItemComponent::SEARCH)) {
      // TODO(omkar): initialize this component
    }
  }

  std::vector<std::future<void>> queue;
  const size_t queue_limit = 10;

  while (!terminating) {
    queue.erase(std::remove_if(queue.begin(), queue.end(),
                               [](std::future<void>& task) {
                                 return task.wait_for(std::chrono::seconds(
                                            0)) == std::future_status::ready;
                               }),
                queue.end());

    if (queue.size() < queue_limit) {
      MediaItemProcessResponse response = api_client->GetMediaItemProcess();
      if (response.id() != "") {
        auto task_future =
            ProcessMediaItem(response, api_client, metadata_component,
                             previewthumbnail_component, places_component,
                             faces_component, ocr_component);
        queue.push_back(std::move(task_future));
      }
    } else {
      SPDLOG_DEBUG("waiting for resources to be free");
    }
  }

  SPDLOG_INFO("waiting for {} active tasks to complete", queue.size());
  for (auto& task : queue) {
    task.wait();
  }
  SPDLOG_INFO("finished all tasks");

  return 0;
}
