import { describe, it, expect } from "vitest";
import { getErrorMessage } from "./client";

describe("getErrorMessage", () => {
  it("uses the error message from the server", () => {
    const err = { response: { data: { error: "email already registered" } } };
    expect(getErrorMessage(err)).toBe("email already registered");
  });

  it("falls back when the server sends no error field", () => {
    const err = { response: { data: {} } };
    expect(getErrorMessage(err, "Failed to load")).toBe("Failed to load");
  });

  it("explains when the server cannot be reached", () => {
    const err = { request: {} };
    expect(getErrorMessage(err)).toMatch(/cannot reach the server/i);
  });
});
