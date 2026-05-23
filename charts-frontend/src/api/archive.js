import { request } from "./client";
import { API_URLS } from "./config";

export function createEpisode(chartId, data) {
  return request(`${API_URLS.charts}/charts/${chartId}/episode`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function getEpisodeById(id) {
  return request(`${API_URLS.archive}/episodes/${id}`);
}

export function getEpisodesByChart(chartId) {
  return request(`${API_URLS.archive}/episodes/chart/${chartId}`);
}

export function getLatestEpisodes(page, pageSize) {
  return request(`${API_URLS.archive}/episodes/latest-page?page=${page}&limit=${pageSize}`);
}

export function searchTracks(query) {
  return request(`${API_URLS.archive}/tracks/search?q=${query}`);
}