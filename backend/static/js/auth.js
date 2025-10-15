
const API = "/api";
export async function apiLogin(email, password, remember) {
  const res = await fetch(`${API}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password })
  });
  if (!res.ok) {
    let msg = "Invalid credentials";
    try { const j = await res.json(); msg = j.error || msg } catch {}
    throw new Error(msg);
  }
  const { accessToken, expiresIn } = await res.json();
  const storage = remember ? localStorage : sessionStorage;
  storage.setItem("accessToken", accessToken);
  storage.setItem("accessTokenExp", (Date.now() + expiresIn * 1000).toString());
  return true;
}