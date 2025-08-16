import { useEffect, useState } from "react";
import { createUser, getUsers, updateUser } from "../../api/admin";

export function meta() {
  return [{ title: "Admin Users - Smriti" }];
}

function CreateUser() {
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleCreateUser = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      let res = await createUser({ name, username, password });
      setSuccess(res);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleCreateUser}>
      {success && (
        <div style={{ color: "green" }}>{JSON.stringify(success)}</div>
      )}
      {error && <div style={{ color: "red" }}>{JSON.stringify(error)}</div>}
      <input
        type="text"
        placeholder="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? "Creating User..." : "Create User"}
      </button>
    </form>
  );
}

function UpdateUser(user) {
  const [id, _] = useState(user.id);
  const [name, setName] = useState(user.name);
  const [username, setUsername] = useState(user.username);
  const [password, setPassword] = useState(user.password);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleUpdateUser = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setSuccess("");

    try {
      let res = await updateUser(id, { name, username, password });
      setSuccess(res);
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleUpdateUser}>
      {success && (
        <div style={{ color: "green" }}>{JSON.stringify(success)}</div>
      )}
      {error && <div style={{ color: "red" }}>{JSON.stringify(error)}</div>}
      <input
        type="text"
        placeholder="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
        required
      />
      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        required
      />
      <button type="submit" disabled={loading}>
        {loading ? "Updating User..." : "Update User"}
      </button>
    </form>
  );
}

function ListUsers() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    try {
      setLoading(true);
      let res = getUsers();
      res.then((data) => {
        setUsers(data);
      });
    } catch (err) {
      setError(err);
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <div>
      <p>Users</p>
      {(!loading &&
        users &&
        users.length > 0 &&
        users.map((user) => <UpdateUser key={user.id} {...user} />)) || (
        <div>
          {error && <div style={{ color: "red" }}>{JSON.stringify(error)}</div>}
        </div>
      )}
    </div>
  );
}

export default function Users() {
  return (
    <div>
      <h1>Admin Users</h1>
      <CreateUser />
      <ListUsers />
    </div>
  );
}
