import { request } from "./client";
import { API_URLS } from "./config";

export function createChart(data) {
  return request(`${API_URLS.charts}/charts/`, {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function getMyChart() {
  return request(`${API_URLS.charts}/charts/me`);
}

export function updateChart(id, data) {
  return request(`${API_URLS.charts}/charts/${id}`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}

export function getChartById(id) {
  return request(`${API_URLS.charts}/charts/${id}`);
}

export function getChartByIdWithoutView(id) {
  return request(`${API_URLS.charts}/charts/guest/${id}`);
}

export function getMyLikedCharts() {
  return request(`${API_URLS.charts}/charts/me/likes`);
}

export function getMyDislikedCharts() {
  return request(`${API_URLS.charts}/charts/me/dislikes`);
}

export function getPopularCharts() {
  return request(`${API_URLS.charts}/charts/popular`);
}