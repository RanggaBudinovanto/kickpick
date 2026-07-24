import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { FormField } from "./form-field";

describe("FormField", () => {
  it("renders the label and forwards input props", () => {
    render(<FormField label="Email" type="email" placeholder="you@example.com" />);

    expect(screen.getByText("Email")).toBeInTheDocument();
    const input = screen.getByPlaceholderText("you@example.com");
    expect(input).toHaveAttribute("type", "email");
  });

  it("does not show an error message when none is passed", () => {
    render(<FormField label="Email" />);
    expect(screen.queryByRole("paragraph")).not.toBeInTheDocument();
  });

  it("shows the error message and marks the input invalid", () => {
    render(<FormField label="Password" error="Password minimal 8 karakter" />);

    expect(screen.getByText("Password minimal 8 karakter")).toBeInTheDocument();
    expect(screen.getByLabelText("Password")).toHaveAttribute("aria-invalid", "true");
  });
});
