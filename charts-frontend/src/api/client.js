import { removeToken } from "../utils/auth";
let sessionExpiredHandled = false;
function getAuthHeaders() {
  const token = localStorage.getItem("token");

  return token
    ? {
        Authorization: `Bearer ${token}`,
      }
    : {};
}

async function request(url, options = {}) {
  const res = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaders(),
      ...(options.headers || {}),
    },
  });
  if (
    res.status === 401 &&
    !url.includes("/auth/login") &&
    !url.includes("/auth/register")
  ) {
    if (sessionExpiredHandled) {
      return null;
    }

    sessionExpiredHandled = true;

    removeToken();

    setTimeout(() => {
      alert("Вы не авторизованы или ваша сессия истекла. Войдите в систему.");
      window.location.href = "/";
    }, 50);

    return null;
  }
  if (res.status === 204) {
    return null;
  }

  if (!res.ok) {
    throw new Error("Ошибка запроса");
  }

  return res.json();
}
export {request};