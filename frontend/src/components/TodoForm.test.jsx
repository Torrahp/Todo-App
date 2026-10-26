import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TodoForm from "./TodoForm";

describe("TodoForm", () => {
  it("disables Add while the input is blank", async () => {
    render(<TodoForm onSubmit={vi.fn()} />);
    expect(screen.getByRole("button", { name: /add/i })).toBeDisabled();

    await userEvent.type(screen.getByLabelText("New todo"), "   ");
    expect(screen.getByRole("button", { name: /add/i })).toBeDisabled();
  });

  it("submits the trimmed title and clears the input on success", async () => {
    const onSubmit = vi.fn().mockResolvedValue(true);
    render(<TodoForm onSubmit={onSubmit} />);
    const input = screen.getByLabelText("New todo");

    await userEvent.type(input, "  buy milk  ");
    await userEvent.click(screen.getByRole("button", { name: /add/i }));

    expect(onSubmit).toHaveBeenCalledWith("buy milk");
    expect(input).toHaveValue("");
  });

  it("keeps the text when saving fails", async () => {
    const onSubmit = vi.fn().mockResolvedValue(false);
    render(<TodoForm onSubmit={onSubmit} />);
    const input = screen.getByLabelText("New todo");

    await userEvent.type(input, "buy milk");
    await userEvent.click(screen.getByRole("button", { name: /add/i }));

    expect(input).toHaveValue("buy milk");
  });
});
