import { createContext } from "react";

export const AuthContext = createContext(null);

// ช่องส่งข้อมูล ให้ component อื่นๆ ที่อยู่ใน <AuthProvider> สามารถเข้าถึงได้