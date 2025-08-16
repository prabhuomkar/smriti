import { api } from "../api";

const login = async (username, password) => {
  const basicAuth = btoa(`${username}:${password}`);
  localStorage.setItem("adminToken", basicAuth);
  const response = await api.get("/v1/users", {
    headers: {
      Authorization: `Basic ${basicAuth}`,
    },
  });

  if (response.status !== 200) {
    localStorage.removeItem("adminToken");
    throw new Error("failed to login");
  }
};

const getUsers = async () => {
  const response = await api.get("/v1/users");
  if (response.status === 200) {
    return response.data;
  }
  throw new Error("failed to get users");
};

const createUser = async (user) => {
  const response = await api.post("/v1/users", user);
  if (response.status === 201) {
    return response.data;
  }
  throw new Error("failed to create user");
};

const updateUser = async (id, user) => {
  const response = await api.put(`/v1/users/${id}`, user);
  if (response.status === 204) {
    return response.data;
  }
  throw new Error("failed to update user");
};

const deleteUser = async (id) => {
  const response = await api.delete(`/v1/users/${id}`);
  if (response.status === 204) {
    return response.data;
  }
  throw new Error("failed to delete user");
};

export { login, getUsers, createUser, updateUser, deleteUser };
