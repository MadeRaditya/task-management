import axios from "axios";

const baseURL = process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080";

export const api = axios.create({
  baseURL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});

api.interceptors.response.use(
  (res) => res,
  (error) => {
    const apiError = error?.response?.data;
    if (apiError?.error) {
      return Promise.reject({
        code: apiError.error.code,
        message: apiError.error.message,
        status: error.response.status,
      });
    }
    return Promise.reject({
      code: "NETWORK_ERROR",
      message: error?.message ?? "Network error",
      status: 0,
    });
  },
);
