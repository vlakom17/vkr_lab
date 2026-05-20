import { request } from "./client";
import { API_URLS } from "./config";

export function sendChartReaction(chartId, type) {
  return request(`${API_URLS.charts}/charts/${chartId}/reaction`, {
    method: "POST",
    body: JSON.stringify({ type }),
  });
}

export function getChartStats(chartId) {
  return request(`${API_URLS.analytics}/analysis/${chartId}`);
}

export function getRecommendations() {
  return request(`${API_URLS.analytics}/analysis/recommendations`);
}

export function getMyReaction(chartId) {
  return request(`${API_URLS.analytics}/analysis/${chartId}/my-reaction`);
}