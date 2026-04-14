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

  if (res.status === 401) {
    removeToken();
    window.location.href = "/login";
    throw new Error("Unauthorized");
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