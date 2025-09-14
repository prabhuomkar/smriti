// Copyright 2025 Omkar Prabhu
#include "worker/metadata.h"

#include <grpcpp/grpcpp.h>
#include <spdlog/spdlog.h>

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

std::unordered_map<std::string, std::string> Metadata::Extract(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& file_path) {
  std::unordered_map<std::string, std::string> result;

  // default values
  result["status"] = "PROCESSING";
  result["type"] = "unknown";
  result["category"] = "default";
  result["exifdata"] = "";

  std::unordered_map<std::string, std::string> exif_data =
      exif_tool_client_->Extract(file_path);

  for (const auto& [key, value] : exif_data) {
    spdlog::debug("exiftool response: {}={}", key, value);
  }

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

  // other fields
  result["mime_type"] = GetValue(exif_data, {"MIMEType"});
  if (result["mime_type"].find("image") != std::string::npos) {
    result["type"] = "photo";
  } else if (result["mime_type"].find("video") != std::string::npos) {
    result["type"] = "video";
  }
  result["width"] =
      GetValue(exif_data, {"ExifImageWidth", "ImageWidth", "SourceImageWidth"});
  result["height"] = GetValue(
      exif_data, {"ExifImageHeight", "ImageHeight", "SourceImageHeight"});
  result["creation_time"] =
      GetValue(exif_data, {"CreateDate", "DateCreated", "DateTimeOriginal"});
  std::regex date_re(R"((\d{4}):(\d{2}):(\d{2}) (\d{2}):(\d{2}):(\d{2}))");
  std::smatch date_match;
  if (std::regex_search(result["creation_time"], date_match, date_re)) {
    result["creation_time"] = date_match[1].str() + "-" + date_match[2].str() +
                              "-" + date_match[3].str() + " " +
                              date_match[4].str() + ":" + date_match[5].str() +
                              ":" + date_match[6].str();
  }

  // camera
  result["camera_make"] = GetValue(exif_data, {"Make", "LensMake"});
  result["camera_model"] = GetValue(exif_data, {"Model", "LensModel"});
  result["focal_length"] = GetValue(exif_data, {"FocalLength"});
  result["aperture_fnumber"] = GetValue(exif_data, {"FNumber"});
  result["iso_equivalent"] = GetValue(exif_data, {"ISO"});
  result["exposure_time"] = GetValue(exif_data, {"ExposureTime"});
  result["fps"] = GetValue(exif_data, {"VideoFrameRate"});
  result["megapixels"] = GetValue(exif_data, {"Megapixels"});

  // gps
  result["longitude"] = GetCoordinates(GetValue(exif_data, {"GPSLongitude"}));
  result["latitude"] = GetCoordinates(GetValue(exif_data, {"GPSLatitude"}));

  // category
  if (exif_data.find("UserComment") != exif_data.end()) {
    std::string value = exif_data["UserComment"];
    std::transform(value.begin(), value.end(), value.begin(),
                   [](unsigned char c) { return std::tolower(c); });
    if (value == "screenshot") {
      result["category"] = "screenshot";
    }
  } else if (exif_data.find("LivePhotoVideoIndex") != exif_data.end()) {
    result["category"] = "live";
  } else if (exif_data.find("FullFrameRatePlaybackIntent") != exif_data.end()) {
    if (exif_data["FullFrameRatePlaybackIntent"] == "0") {
      result["category"] = "timelapse";
    }
  } else if (exif_data.find("CaptureMode") != exif_data.end()) {
    std::string value = exif_data["CaptureMode"];
    std::transform(value.begin(), value.end(), value.begin(),
                   [](unsigned char c) { return std::tolower(c); });
    if (value == "time-lapse") {
      result["category"] = "timelapse";
    }
  } else if (result["width"] != "" && result["height"] != "") {
    int width = std::stoi(result["width"]);
    int height = std::stoi(result["height"]);
    if (width >= 10000 && (height * 4 <= width)) {
      result["category"] = "panorama";
    }
  } else if (result["fps"] != "") {
    int fps = std::stoi(result["fps"]);
    if (fps >= 120) {
      result["category"] = "slow";
    }
  }

  MediaItemMetadataRequest request;
  request.set_userid(user_id);
  request.set_mediaitemid(mediaitem_id);
  request.set_status(result["status"]);
  request.set_mimetype(result["mime_type"]);
  request.set_type(result["type"]);
  request.set_category(result["category"]);
  if (result["width"] != "") {
    request.set_width(std::stoi(result["width"]));
  }
  if (result["height"] != "") {
    request.set_height(std::stoi(result["height"]));
  }
  request.set_creationtime(result["creation_time"]);
  request.set_cameramake(result["camera_make"]);
  request.set_cameramodel(result["camera_model"]);
  request.set_focallength(result["focal_length"]);
  request.set_aperturefnumber(result["aperture_fnumber"]);
  request.set_isoequivalent(result["iso_equivalent"]);
  request.set_exposuretime(result["exposure_time"]);
  request.set_megapixels(result["megapixels"]);
  request.set_fps(result["fps"]);
  if (result["latitude"] != "") {
    request.set_latitude(std::stod(result["latitude"]));
  }
  if (result["longitude"] != "") {
    request.set_longitude(std::stod(result["longitude"]));
  }
  request.set_exifdata(result["exifdata"]);
  bool ok = api_client_->SaveMediaItemMetadata(request);

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
