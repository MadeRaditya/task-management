import { useCallback, useEffect, useState } from "react";
import { listTasks } from "../api/tasks";
import { Task, TaskStatus } from "../types/task";

interface UseTasksState {
  tasks: Task[];
  loading: boolean;
  error: string | null;
  total: number;
  totalPages: number;
}

export interface UseTasksFilters {
  keyword: string;
  status: TaskStatus | "";
  page: number;
  limit: number;
}

export function useTasks(initial?: Partial<UseTasksFilters>) {
  const [filters, setFilters] = useState<UseTasksFilters>({
    keyword: initial?.keyword ?? "",
    status: initial?.status ?? "",
    page: initial?.page ?? 1,
    limit: initial?.limit ?? 10,
  });

  const [state, setState] = useState<UseTasksState>({
    tasks: [],
    loading: false,
    error: null,
    total: 0,
    totalPages: 0,
  });

  const fetchTasks = useCallback(async () => {
    setState((s) => ({ ...s, loading: true, error: null }));
    try {
      const result = await listTasks({
        keyword: filters.keyword || undefined,
        status: filters.status || undefined,
        page: filters.page,
        limit: filters.limit,
        sort_by: "created_at",
        sort_dir: "desc",
      });
      setState({
        tasks: result.items,
        loading: false,
        error: null,
        total: result.total,
        totalPages: result.total_pages,
      });
    } catch (e: any) {
      setState((s) => ({
        ...s,
        loading: false,
        error: e?.message ?? "Gagal memuat data",
      }));
    }
  }, [filters.keyword, filters.status, filters.page, filters.limit]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  const setKeyword = (keyword: string) =>
    setFilters((f) => ({ ...f, keyword, page: 1 }));
  const setStatus = (status: TaskStatus | "") =>
    setFilters((f) => ({ ...f, status, page: 1 }));
  const setPage = (page: number) => setFilters((f) => ({ ...f, page }));

  return {
    ...state,
    filters,
    setKeyword,
    setStatus,
    setPage,
    refresh: fetchTasks,
  };
}
