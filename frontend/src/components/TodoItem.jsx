import { useState } from "react";
import { Check, Pencil, Trash2, X } from "lucide-react";

function TodoItem({ todo, onToggle, onRename, onDelete }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(todo.title);
  const [toggling, setToggling] = useState(false);

  const handleToggle = async () => {
    setToggling(true);
    await onToggle(todo);
    setToggling(false);
  };

  const startEdit = () => {
    setDraft(todo.title);
    //เอาค่าtitle มาใส่
    setEditing(true);
  };

  const handleSave = async (e) => {
    e.preventDefault();// ป้องกันการ rerender 
    const trimmed = draft.trim();
    if (!trimmed) return;
    if (trimmed === todo.title) {
        //ถ้าชื่อเดี่ยวกันก็แค่ ปิด edit
      setEditing(false);
      return;
    }
    const ok = await onRename(todo.id, trimmed);
    if (ok) setEditing(false);
  };

  const handleDelete = () => {
    if (window.confirm(`Delete "${todo.title}"?`)) onDelete(todo.id);
  };

  if (editing) {
    return (
      <li className="todo-item">
        <form className="todo-edit" onSubmit={handleSave}>
          <input
            type="text"
            aria-label="Edit todo title"
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Escape") setEditing(false);
            }}
            maxLength={255}
            autoFocus
          />
          <button type="submit">
            <Check size={18} strokeWidth={1.75} aria-hidden="true" />
            Save
          </button>
          <button type="button" className="ghost" onClick={() => setEditing(false)}>
            <X size={16} strokeWidth={1.75} aria-hidden="true" />
            Cancel
          </button>
        </form>
      </li>
    );
  }

  return (
    <li className={todo.is_done ? "todo-item done" : "todo-item"}>
      <label className="todo-check">
        <input
          type="checkbox"
          checked={todo.is_done}
          onChange={handleToggle}
          disabled={toggling}
        />
        <span>{todo.title}</span>
      </label>
      <div className="todo-actions">
        <button
          type="button"
          className="ghost icon-btn"
          aria-label="Edit"
          title="Edit"
          onClick={startEdit}
        >
          <Pencil size={18} strokeWidth={1.75} aria-hidden="true" />
        </button>
        <button
          type="button"
          className="ghost icon-btn danger"
          aria-label="Delete"
          title="Delete"
          onClick={handleDelete}
        >
          <Trash2 size={18} strokeWidth={1.75} aria-hidden="true" />
        </button>
      </div>
    </li>
  );
}

export default TodoItem;
