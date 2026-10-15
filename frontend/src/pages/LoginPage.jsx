import { useState } from "react";
import { Navigate } from "react-router-dom";
import { CircleAlert, Lock, LogIn, Mail, UserPlus } from "lucide-react";
import { useAuth } from "../hooks/useAuth";
import { getErrorMessage } from "../api/client";

function LoginPage() {
  const { isAuthenticated, login, register } = useAuth();
  const [mode, setMode] = useState("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);
  const [submitting, setSubmitting] = useState(false);

  if (isAuthenticated) return <Navigate to="/" replace />;

  const isRegister = mode === "register";//จะเปิด True ก็ต่อเมื่อ mode เป็น register
  const submitLabel = isRegister ? "Create account" : "Log in";
  const SubmitIcon = isRegister ? UserPlus : LogIn;

  const handleSubmit = async (e) => {
    e.preventDefault();// ยกเลิกพฤติกรรมของbrowser เช่นการ reload หน้า
    setError(null);
    setSubmitting(true);
    try {
      if (isRegister) {
        await register(email, password);
      } else {
        await login(email, password);
      }
    } catch (err) {
      setError(getErrorMessage(err, "Failed to sign in"));
    } finally {
      setSubmitting(false);
    }
  };

  const switchMode = () => {
    setMode(isRegister ? "login" : "register");
    setError(null);
  };

  return (
    <main className="page narrow">
      <h1>{submitLabel}</h1>

      <section className="sheet">
        <form className="auth-form" onSubmit={handleSubmit}>
          <label>
            Email
            <span className="field">
              <Mail size={18} strokeWidth={1.75} aria-hidden="true" />
              <input
                type="email"
                autoComplete="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </span>
          </label>
          <label>
            Password
            <span className="field">
              <Lock size={18} strokeWidth={1.75} aria-hidden="true" />
              <input
                type="password"
                autoComplete={isRegister ? "new-password" : "current-password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                minLength={8}
                required
              />
            </span>
          </label>
          {error && (
            <p role="alert" className="error">
              <CircleAlert size={18} strokeWidth={1.75} aria-hidden="true" />
              <span>{error}</span>
            </p>
          )}
          <button type="submit" disabled={submitting}>
            <SubmitIcon size={18} strokeWidth={1.75} aria-hidden="true" />
            {submitting ? "Please wait…" : submitLabel}
          </button>
        </form>
      </section>

      <p className="muted">
        {isRegister ? "Already have an account? " : "No account yet? "}
        <button type="button" className="link" onClick={switchMode}>
          {isRegister ? "Log in" : "Create one"}
        </button>
      </p>
    </main>
  );
}

export default LoginPage;
