// Copyright 2025 Omkar Prabhu
#include "worker/metadata.h"

#include <grpcpp/grpcpp.h>
#include <spdlog/spdlog.h>
#include <sys/wait.h>

#include <algorithm>
#include <cstdio>
#include <memory>
#include <optional>
#include <regex>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "protos/api.pb.h"
#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace metadata {

ExifToolClient::ExifToolClient() {
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

ExifToolClient::~ExifToolClient() {
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

std::unordered_map<std::string, std::string> ExifToolClient::Extract(
    const std::string& file_path) {
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

std::unordered_map<std::string, std::string> Metadata::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  std::unordered_map<std::string, std::string> result;
  std::unordered_map<std::string, std::string> exif_data;

  // default values
  MediaItemStatus status = MediaItemStatus::PROCESSING;
  MediaItemType type = MediaItemType::UNKNOWN;
  MediaItemCategory category = MediaItemCategory::DEFAULT;

  MediaItemMetadataRequest request;
  request.set_userid(user_id);
  request.set_mediaitemid(mediaitem_id);

  try {
    exif_data = exif_tool_client_->Extract(file_path);

    SPDLOG_DEBUG("exiftool response: {}", exif_data);

    std::string raw_exifdata = "{";
    for (auto it = exif_data.begin(); it != exif_data.end();) {
      const std::string& key = it->first;
      std::string value = it->second;

      // raw exif data
      raw_exifdata.append("\"");
      raw_exifdata.append(key);
      raw_exifdata.append("\":\"");
      raw_exifdata.append(value);
      raw_exifdata.append("\"");
      if (++it != exif_data.end()) {
        raw_exifdata.append(",");
      }
    }
    raw_exifdata.append("}");
    result["exifdata"] = raw_exifdata;
    request.set_exifdata(result["exifdata"]);

    // other fields
    result["mime_type"] = GetValue(exif_data, {"MIMEType"});
    request.set_mimetype(result["mime_type"]);
    if (result["mime_type"].find("image") != std::string::npos) {
      type = MediaItemType::PHOTO;
    } else if (result["mime_type"].find("video") != std::string::npos) {
      type = MediaItemType::VIDEO;
    }
    result["width"] = GetValue(
        exif_data, {"ExifImageWidth", "ImageWidth", "SourceImageWidth"});
    result["height"] = GetValue(
        exif_data, {"ExifImageHeight", "ImageHeight", "SourceImageHeight"});
    if (result["width"] != "") {
      request.set_width(std::stoi(result["width"]));
    }
    if (result["height"] != "") {
      request.set_height(std::stoi(result["height"]));
    }
    result["creation_time"] =
        GetValue(exif_data, {"CreateDate", "DateCreated", "DateTimeOriginal"});
    std::regex date_re(R"((\d{4}):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2}))");
    std::smatch date_match;
    if (std::regex_search(result["creation_time"], date_match, date_re)) {
      result["creation_time"] =
          date_match[1].str() + "-" + date_match[2].str() + "-" +
          date_match[3].str() + " " + date_match[4].str() + ":" +
          date_match[5].str() + ":" + date_match[6].str();
    }

    // camera
    result["camera_make"] = GetValue(exif_data, {"Make", "LensMake"});
    request.set_cameramake(result["camera_make"]);

    result["camera_model"] = GetValue(exif_data, {"Model", "LensModel"});
    request.set_cameramodel(result["camera_model"]);

    result["focal_length"] = GetValue(exif_data, {"FocalLength"});
    request.set_focallength(result["focal_length"]);

    result["aperture_fnumber"] = GetValue(exif_data, {"FNumber"});
    request.set_aperturefnumber(result["aperture_fnumber"]);

    result["iso_equivalent"] = GetValue(exif_data, {"ISO"});
    request.set_isoequivalent(result["iso_equivalent"]);

    result["exposure_time"] = GetValue(exif_data, {"ExposureTime"});
    request.set_exposuretime(result["exposure_time"]);

    result["fps"] = GetValue(exif_data, {"VideoFrameRate"});
    request.set_fps(result["fps"]);

    result["megapixels"] = GetValue(exif_data, {"Megapixels"});
    request.set_megapixels(result["megapixels"]);

    // gps
    result["longitude"] = GetValue(exif_data, {"GPSLongitude"});
    if (result["longitude"] != "") {
      result["longitude"] = GetCoordinates(result["longitude"]);
      request.set_longitude(std::stod(result["longitude"]));
    }
    result["latitude"] = GetValue(exif_data, {"GPSLatitude"});
    if (result["latitude"] != "") {
      result["latitude"] = GetCoordinates(result["latitude"]);
      request.set_latitude(std::stod(result["latitude"]));
    }

    // category
    if (exif_data.find("UserComment") != exif_data.end()) {
      std::string value = exif_data["UserComment"];
      std::transform(value.begin(), value.end(), value.begin(),
                     [](unsigned char c) { return std::tolower(c); });
      if (value == "screenshot") {
        category = MediaItemCategory::SCREENSHOT;
      }
    } else if (exif_data.find("LivePhotoVideoIndex") != exif_data.end()) {
      category = MediaItemCategory::LIVE;
    } else if (exif_data.find("FullFrameRatePlaybackIntent") !=
               exif_data.end()) {
      if (exif_data["FullFrameRatePlaybackIntent"] == "0") {
        category = MediaItemCategory::TIMELAPSE;
      }
    } else if (exif_data.find("CaptureMode") != exif_data.end()) {
      std::string value = exif_data["CaptureMode"];
      std::transform(value.begin(), value.end(), value.begin(),
                     [](unsigned char c) { return std::tolower(c); });
      if (value == "time-lapse") {
        category = MediaItemCategory::TIMELAPSE;
      }
    } else if (result["width"] != "" && result["height"] != "") {
      int width = std::stoi(result["width"]);
      int height = std::stoi(result["height"]);
      if (width >= 10000 && (height * 4 <= width)) {
        category = MediaItemCategory::PANORAMA;
      }
    } else if (result["fps"] != "") {
      int fps = std::stoi(result["fps"]);
      if (fps >= 120) {
        category = MediaItemCategory::SLOW;
      }
    }
  } catch (const std::exception& e) {
    SPDLOG_ERROR("error extracting metadata: {}", e.what());
    status = MediaItemStatus::FAILED;
  }

  result["status"] = MediaItemStatus_Name(status);
  result["type"] = MediaItemType_Name(type);
  result["category"] = MediaItemCategory_Name(category);

  request.set_status(status);
  request.set_type(type);
  request.set_category(category);
  api_client_->SaveMediaItemMetadata(request);

  return result;
}

std::string GetCoordinates(const std::string& location) {
  int degrees = 0, minutes = 0;
  double seconds = 0.0;
  char direction = 'N';
  double decimal = 0.0;

  if (std::sscanf(location.c_str(), "%d%*[^0-9]%d%*[^0-9]%lf%*[^NSEW]%c",
                  &degrees, &minutes, &seconds, &direction) == 4) {
    decimal = degrees + minutes / 60.0 + seconds / 3600.0;
    if (direction == 'S' || direction == 'W') {
      decimal = -decimal;
    }
  }
  return std::to_string(decimal);
}

std::string GetValue(const std::unordered_map<std::string, std::string>& data,
                     const std::vector<std::string>& keys) {
  for (const std::string& key : keys) {
    auto it = data.find(key);
    if (it != data.end()) {
      return it->second; // no double lookup
    }
  }
  return "";
}

std::shared_ptr<Metadata> Init(const ComponentConfig& config,
                               std::shared_ptr<APIClient> api_client) {
  return std::make_shared<Metadata>(std::make_shared<ExifToolClient>(),
                                    api_client);
}

} // namespace metadata

} // namespace components
