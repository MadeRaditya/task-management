import { api } from "./client";
import {
  CreateTaskPayload,
  ListTasksQuery,
  ListTasksResult,
  Task,
  UpdateTaskPayload,
} from "../types/task";

interface DataResponse<T> {
  data: T;
}

export async function listTasks(
  q: ListTasksQuery = {},
): Promise<ListTasksResult> {
  const params: Record<string, string | number> = {};
  if (q.status) params.status = q.status;
  if (q.keyword) params.keyword = q.keyword;
  if (q.assignee) params.assignee = q.assignee;
  if (q.page) params.page = q.page;
  if (q.limit) params.limit = q.limit;
  if (q.sort_by) params.sort_by = q.sort_by;
  if (q.sort_dir) params.sort_dir = q.sort_dir;

  const res = await api.get<ListTasksResult>("/api/tasks", { params });
  return res.data;
}

export async function getTask(id: number): Promise<Task> {
  const res = await api.get<DataResponse<Task>>(`/api/tasks/${id}`);
  return res.data.data;
}

export async function createTask(payload: CreateTaskPayload): Promise<Task> {
  const res = await api.post<DataResponse<Task>>("/api/tasks", payload);
  return res.data.data;
}

export async function updateTask(
  id: number,
  payload: UpdateTaskPayload,
): Promise<Task> {
  const res = await api.put<DataResponse<Task>>(`/api/tasks/${id}`, payload);
  return res.data.data;
}

export async function deleteTask(id: number): Promise<void> {
  await api.delete(`/api/tasks/${id}`);
}
