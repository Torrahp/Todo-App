import { useState } from "react";
import { AuthContext } from "./AuthContext";
import * as authApi from "../api/auth";

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem("token"));
  const [email, setEmail] = useState(() => localStorage.getItem("email"));

  const login = async (userEmail, password) => {
    const data = await authApi.login(userEmail, password);
    localStorage.setItem("token", data.token);
    localStorage.setItem("email", userEmail);
    setToken(data.token);
    setEmail(userEmail);
  };

  const register = async (userEmail, password) => {
    await authApi.register(userEmail, password);
    await login(userEmail, password);
  };

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("email");
    setToken(null);
    setEmail(null);
  };

  const value = {
    token,
    email,
    isAuthenticated: Boolean(token),//แค่ตรวจสอบว่ามี token หรือไม่ ถ้ามีคือ true ถ้าไม่มีคือ false
    login,
    register,
    logout,
  };

  return <AuthContext value={value}>{children}</AuthContext>;
}
