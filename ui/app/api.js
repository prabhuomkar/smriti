import axios from "axios";

const BASIC_AUTH_ROUTES = ["/v1/users", "/v1/jobs"];

export const api = axios.create({
  baseURL: "http://localhost:5001",
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.request.use(
  (config) => {
    if (BASIC_AUTH_ROUTES.find((route) => config.url.includes(route))) {
      config.headers["Authorization"] =
        `Basic ${localStorage.getItem("adminToken")}`;
    } else {
      config.headers["Authorization"] =
        `Bearer ${localStorage.getItem("accessToken")}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export const getErrorMessage = (status) => {
  switch (status) {
    case 401:
      return "Incorrect username or password";
    case 400:
      return "Bad Request";
    default:
      return "Some error. Try again!";
  }
};
