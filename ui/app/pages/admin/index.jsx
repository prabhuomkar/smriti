import { useState } from "react";
import { login } from "../../api/admin";
import { useNavigate } from "react-router";
import { getErrorMessage } from "../../api";

export function meta() {
  return [{ title: "Admin Login - Smriti" }];
}

export default function Login() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await login(username, password);
      navigate("/admin/users");
    } catch (err) {
      setError(getErrorMessage(err?.response?.status));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="w-full max-w-sm space-y-4 p-4">
        <div className="flex flex-row items-center justify-between">
          <div className="smriti-dark text-lg font-medium align-center">
            Admin Login
          </div>
          <div className="text-sm font-medium align-center">
            or{" "}
            <a href="https://smriti.omkar.xyz/docs/user-guide/deployment">
              setup an admin
            </a>
          </div>
        </div>
        <form className="flex flex-col space-y-4" onSubmit={handleLogin}>
          <input
            name="username"
            className="py-2 px-3 border border-gray-300 outline-none rounded-md w-full"
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => {
              setUsername(e.target.value);
              setError("");
            }}
            autoComplete="off"
            required
          />
          <input
            name="password"
            className="py-2 px-3 border border-gray-300 outline-none rounded-md"
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => {
              setPassword(e.target.value);
              setError("");
            }}
            autoComplete="off"
            required
          />
          {error && (
            <div className="smriti-error smriti-error-bg border smriti-error-border p-2 rounded-md flex items-center space-x-2 justify-center">
              <span>{error}</span>
            </div>
          )}
          <button
            className="smriti-bg text-white p-2 mt-4 transition rounded-md"
            type="submit"
            disabled={loading}
          >
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>
      </div>
    </div>
  );
}
