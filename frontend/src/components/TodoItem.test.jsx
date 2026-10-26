import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TodoItem from "./TodoItem";

const todo = { id: 7, title: "buy milk", is_done: false };

function renderItem(overrides = {}) {
  const props = {
    todo,
    onToggle: vi.fn(),
    onRename: vi.fn().mockResolvedValue(true),
    onDelete: vi.fn(),
    ...overrides,
  };
  render(<TodoItem {...props} />);
  return props;
}

describe("TodoItem", () => {
  it("calls onToggle with the todo when the checkbox is clicked", async () => {
    const { onToggle } = renderItem();
    await userEvent.click(screen.getByRole("checkbox"));
    expect(onToggle).toHaveBeenCalledWith(todo);
  });

  it("renames through the edit form", async () => {
    const { onRename } = renderItem();
    await userEvent.click(screen.getByRole("button", { name: "Edit" }));

    const input = screen.getByLabelText("Edit todo title");
    await userEvent.clear(input);
    await userEvent.type(input, "buy oat milk");
    await userEvent.click(screen.getByRole("button", { name: /save/i }));

    expect(onRename).toHaveBeenCalledWith(7, "buy oat milk");
  });

  it("cancels editing with Escape", async () => {
    const { onRename } = renderItem();
    await userEvent.click(screen.getByRole("button", { name: "Edit" }));
    await userEvent.keyboard("{Escape}");

    expect(screen.queryByLabelText("Edit todo title")).not.toBeInTheDocument();
    expect(onRename).not.toHaveBeenCalled();
  });

  it("deletes only after the user confirms", async () => {
    const confirmSpy = vi.spyOn(window, "confirm");
    const { onDelete } = renderItem();

    confirmSpy.mockReturnValueOnce(false);
    await userEvent.click(screen.getByRole("button", { name: "Delete" }));
    expect(onDelete).not.toHaveBeenCalled();

    confirmSpy.mockReturnValueOnce(true);
    await userEvent.click(screen.getByRole("button", { name: "Delete" }));
    expect(onDelete).toHaveBeenCalledWith(7);

    confirmSpy.mockRestore();
  });
});
