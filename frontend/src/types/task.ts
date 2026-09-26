export type TaskStatus = "todo" | "in_progress" | "done";

export interface Task {
  id: number;
  title: string;
  description?: string;
  status: TaskStatus;
  assignee?: string;
  created_at: string;
  updated_at: string;
}

export interface ListTasksQuery {
  status?: TaskStatus | "";
  keyword?: string;
  assignee?: string;
  page?: number;
  limit?: number;
  sort_by?: string;
  sort_dir?: "asc" | "desc";
}

export interface ListTasksResult {
  items: Task[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface CreateTaskPayload {
  title: string;
  description?: string;
  status?: TaskStatus;
  assignee?: string;
}

export interface UpdateTaskPayload {
  title?: string;
  description?: string;
  status?: TaskStatus;
  assignee?: string;
}

export interface ApiError {
  error: {
    code: string;
    message: string;
  };
}
