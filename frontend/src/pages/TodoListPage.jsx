import {
  CircleAlert,
  ClipboardList,
  LoaderCircle,
  LogOut,
  RotateCw,
} from "lucide-react";
import { useAuth } from "../hooks/useAuth";
import { useTodos } from "../hooks/useTodos";
import TodoForm from "../components/TodoFrom";
import TodoItem from "../components/TodoItem";

function TodoListPage() {
  const { email, logout } = useAuth();
  const {
    todos,
    loading,
    error,
    actionError,
    addTodo,
    toggleTodo,
    renameTodo,
    removeTodo,
    refetch,
  } = useTodos();

  const remaining = todos.filter((t) => !t.is_done).length;

  let content;
  if (loading) {
    content = (
      <p className="muted loading">
        <LoaderCircle className="spin" size={16} strokeWidth={1.75} aria-hidden="true" />
        Loading…
      </p>
    );
  } else if (error) {
    content = (
      <div role="alert" className="error">
        <CircleAlert size={18} strokeWidth={1.75} aria-hidden="true" />
        <span>{error}</span>
        <button type="button" className="ghost" onClick={refetch}>
          <RotateCw size={16} strokeWidth={1.75} aria-hidden="true" />
          Try again
        </button>
      </div>
    );
  } else if (todos.length === 0) {
    content = (
      <div className="muted empty">
        <ClipboardList size={28} strokeWidth={1.5} aria-hidden="true" />
        <p>No todos yet. Add your first one above.</p>
      </div>
    );
  } else {
    content = (
      <>
        <p className="muted">
          {remaining} of {todos.length} remaining
        </p>
        <ul className="todo-list">
          {todos.map((t) => (
            <TodoItem
              key={t.id}
              todo={t}
              onToggle={toggleTodo}
              onRename={renameTodo}
              onDelete={removeTodo}
            />
          ))}
        </ul>
      </>
    );
  }

  return (
    <main className="page">
      <header className="page-header">
        <h1>My Todos</h1>
        <div className="user-box">
          <span>{email}</span>
          <button type="button" className="ghost" onClick={logout}>
            <LogOut size={16} strokeWidth={1.75} aria-hidden="true" />
            Log out
          </button>
        </div>
      </header>

      <section className="sheet">
        <TodoForm onSubmit={addTodo} />

        {actionError && (
          <p role="alert" className="error">
            <CircleAlert size={18} strokeWidth={1.75} aria-hidden="true" />
            <span>{actionError}</span>
          </p>
        )}

        {content}
      </section>
    </main>
  );
}

export default TodoListPage;
