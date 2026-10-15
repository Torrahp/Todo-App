import { useState, useEffect } from "react";
import * as todosApi from "../api/todos";
import { getErrorMessage } from "../api/client";

export function useTodos() {
  const [todos, setTodos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [actionError, setActionError] = useState(null);
  const [reloadKey, setReloadKey] = useState(0);

  useEffect(() => {
    async function load() {
      try {
        const data = await todosApi.getTodos();
        setTodos(data);// เอาข้อมูลมาใส่
        setError(null);
      } catch (err) {
        setError(getErrorMessage(err, "Failed to load todos"));
      } finally {
        setLoading(false);
      }
    }
    load();
  }, [reloadKey]);//[]เป็น dependency array มีไว้เพื่อให้ react เช็ค หากค่าแรกไม่ตรงกับค่าสอง จะทำการรัน effect ใหม่

  const refetch = () => {
    setLoading(true);
    setActionError(null);
    setReloadKey((key) => key + 1);
  };

  const run = async (action, fallback) => {
    setActionError(null);
    try {
      await action();
      return true;
    } catch (err) {
      setActionError(getErrorMessage(err, fallback));
      return false;
    }
  };

  const replaceTodo = (updated) => {
    setTodos((prev) => prev.map((t) => (t.id === updated.id ? updated : t)));
  };

  const addTodo = (title) =>
    run(async () => {
      const created = await todosApi.createTodo(title);
      setTodos((prev) => [created, ...prev]);
    }, "Failed to add todo");

  const toggleTodo = (todo) =>
    run(async () => {
      replaceTodo(await todosApi.updateTodo(todo.id, { is_done: !todo.is_done }));
    }, "Failed to update todo");

  const renameTodo = (id, title) =>
    run(async () => {
      replaceTodo(await todosApi.updateTodo(id, { title }));
    }, "Failed to rename todo");

  const removeTodo = (id) =>
    run(async () => {
      await todosApi.deleteTodo(id);
      setTodos((prev) => prev.filter((t) => t.id !== id));
    }, "Failed to delete todo");

  return {
    todos,
    loading,
    error,
    actionError,
    addTodo,
    toggleTodo,
    renameTodo,
    removeTodo,
    refetch,
  };
}
