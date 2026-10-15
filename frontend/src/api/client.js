import axios from 'axios';
// แปลง Object เป็น JSON, แปลงJSON ตอนรับ, ถ้า status 400/500 จะthrow ให้เอง, แนบ token ให้เองผ่านinterceptors
const client = axios.create({
     baseURL: import.meta.env.VITE_API_URL ?? "http://localhost:8080"
    });

client.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config
});// มันจะเติม token ก่อนไปที่หน้าบ้าน

client.interceptors.response.use(
    (res) => res,
    (err) => {
        const isAuthRequest = err.config?.url?.startsWith("/auth/");//path เริ่มต้นด้วย /auth
        if(err.response?.status === 401 && ! isAuthRequest){
            localStorage.removeItem("token");
            localStorage.removeItem("email");
            window.location.href = "/login" ;//พากลับไปที่loginPage
        }
        return Promise.reject(err);
    }
);// ถ้า

export function getErrorMessage(err, fallback = "Something went wrong") {
    if (err.response) return err.response.data?.error ?? fallback;
    if (err.request) return "Cannot reach the server. Is the backend running?";
    return fallback;
}

export default client;