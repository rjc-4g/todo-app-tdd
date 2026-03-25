import axios from 'axios';
import type { Task } from '../types/task';

const BASE = '/api/v1/tasks';

export const fetchTasks = async (): Promise<Task[]> => {
  const { data } = await axios.get<Task[]>(BASE);
  return data;
};

export const createTask = async (name: string): Promise<Task> => {
  const { data } = await axios.post<Task>(BASE, { name, status: 0 });
  return data;
};

export const updateTaskStatus = async (id: number, status: number): Promise<Task> => {
  const { data } = await axios.patch<Task>(`${BASE}/${id}`, { status });
  return data;
};

export const deleteTask = async (id: number): Promise<Task> => {
  const { data } = await axios.delete<Task>(`${BASE}/${id}`);
  return data;
};
