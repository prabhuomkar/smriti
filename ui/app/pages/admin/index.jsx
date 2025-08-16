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
        <img src="logo.png" alt="Smriti" width="48" className="block mx-auto" />
        <div className="smriti-dark-color text-md text-center">Admin Login</div>
        <form className="flex flex-col space-y-4" onSubmit={handleLogin}>
          <input
            name="username"
            className="p-2 border border-gray-300 text-sm outline-none"
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
            className="p-2 border border-gray-300 text-sm outline-none"
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
            <div className="text-sm text-red-500 bg-red-100 p-2 flex items-center space-x-1 justify-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
                />
              </svg>
              <span>{error}</span>
            </div>
          )}
          <button
            className="smriti-bg-color text-white text-sm p-2 mt-4 transition"
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
