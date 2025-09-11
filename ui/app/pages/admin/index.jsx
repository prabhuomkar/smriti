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
          <div className="smriti-dark-color text-lg font-medium align-center">
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
            <div className="text-red-500 bg-red-100 border border-red-300 p-2 rounded-md flex items-center space-x-2 justify-center">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.6}
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
            className="smriti-bg-color text-white p-2 mt-4 transition rounded-md"
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
