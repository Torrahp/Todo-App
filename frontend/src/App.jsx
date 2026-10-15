import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./context/AuthProvider";
import ProtectedRoute from "./components/ProtectedRoute";
import LoginPage from "./pages/LoginPage";
import TodoListPage from "./pages/TodoListPage";

function App() {
  return (
    // เป็นตัวเชื่อมเพื่อที่จะทำ client-side routing จะได้ไม่ต้องโหลดหน้าใหม่
    // เก็บข้อมูลพวก token และ user info ไว้ใน context เพื่อให้ component อื่นๆ สามารถเข้าถึงได้
    <BrowserRouter> 
      <AuthProvider>
        <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route 
          path="/" 
          element={
          <ProtectedRoute>
            <TodoListPage />
          </ProtectedRoute>
        } 
        />   {/*ProtectedRoute ใช้ตรวจสอบ ว่ามี token ไหม ถ้าไม่ก็ Redirect ไปหน้า login*/}
        <Route path="*" element={<Navigate to="/" replace />} /> {/*รับทุก path ที่ไม่ตรงกับ route อื่นๆ*/}
      </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;