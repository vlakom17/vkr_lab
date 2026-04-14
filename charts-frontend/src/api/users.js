import { request } from "./client";
import { API_URLS } from "./config";

export function registerUser(name, email, password, about) {
  return request(`${API_URLS.users}/auth/register`, {
    method: "POST",
    body: JSON.stringify({ name, email, password, about }),
  });
}

export function loginUser(email, password) {
  return request(`${API_URLS.users}/auth/login`, {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function getMe() {
  return request(`${API_URLS.users}/users/me`);
}

export function updateProfile(data) {
  return request(`${API_URLS.users}/users/me`, {
    method: "PATCH",
    body: JSON.stringify(data),
  });
}

export function logout() {
  return request(`${API_URLS.users}/auth/logout`, {
    method: "POST",
  });
}

export function getUserById(id) {
  return request(`${API_URLS.users}/users/${id}`);
}
