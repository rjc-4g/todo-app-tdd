import { API_BASE_URL } from "./constants";

export const getTasks = () => {
    return fetch(`${API_BASE_URL}/api/v1/tasks`);
};

export const postTask = (name: string) => {
    return fetch(`${API_BASE_URL}/api/v1/tasks`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ name: name }),
    });
};

export const deleteTask = (id: number) => {
    return fetch(`${API_BASE_URL}/api/v1/tasks/${id}`, {
      method: "DELETE",
    });
};
