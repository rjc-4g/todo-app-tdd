import { API_BASE_URL, HEADERS } from "./constants";

const request = async (path: string, options: RequestInit = {}) => {
  try {

    const response = await fetch(`${API_BASE_URL}${path}`, options);

    if (!response.ok) {
      console.error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response;

  } catch (error) {
    console.error("Network Error:", error);
    return undefined;
  }
};

export const getTasks = () => {
    return request("/api/v1/tasks");
};

export const postTask = (name: string) => {
    return request("/api/v1/tasks", {
      method: "POST",
      headers: HEADERS,
      body: JSON.stringify({ name: name }),
    });
};

export const patchTask = (id: number, status: number) => {
    return request(`/api/v1/tasks/${id}`, {
      method: "PATCH",
      headers: HEADERS,
      body: JSON.stringify({ status: status }),
    });
};

export const deleteTask = (id: number) => {
    return request(`/api/v1/tasks/${id}`, {
      method: "DELETE",
    });
};
