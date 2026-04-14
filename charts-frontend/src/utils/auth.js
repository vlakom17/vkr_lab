export function getToken() {
  return localStorage.getItem("token");
}

export function setToken(token) {
  localStorage.setItem("token", token);
  window.dispatchEvent(new Event("authChanged"));
}

export function removeToken() {
  localStorage.removeItem("token");
  window.dispatchEvent(new Event("authChanged"));
}

export function isAuthenticated() {
  return !!getToken();
}