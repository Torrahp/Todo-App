import { useState } from "react";
import { Plus } from "lucide-react";

function TodoForm({ onSubmit }) {
  const [title, setTitle] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();// กันไม่ใช้ reload หน้าเว็ป
    const trimmed = title.trim();
    if (!trimmed) return;// เช็คว่ามีค่าในtrimmed ไหม ถ้าไม่จะreturn ออกเลย
    setSubmitting(true);// บอกว่ากำลังโหลดนะ
    const ok = await onSubmit(trimmed);
    setSubmitting(false);// บอกว่าโหลดแล้วนะ
    if (ok) setTitle("");
  };

  return (
    <form className="todo-form" onSubmit={handleSubmit}>
      <input
        type="text"
        aria-label="New todo"
        placeholder="What needs to be done?"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        maxLength={255}
      />
      <button type="submit" disabled={submitting || !title.trim()}>
        <Plus size={18} strokeWidth={1.75} aria-hidden="true" />
        Add
      </button>
    </form>
  );
}

export default TodoForm;
