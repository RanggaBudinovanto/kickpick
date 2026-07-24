import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";
import messages from "@/messages/id.json";
import LoginPage from "./page";

const pushMock = vi.fn();
const mutateMock = vi.fn();

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
  useRouter: () => ({ push: pushMock }),
}));

vi.mock("@/hooks/use-auth", () => ({
  useLogin: () => ({
    mutate: mutateMock,
    isPending: false,
  }),
}));

function renderLoginPage() {
  return render(
    <NextIntlClientProvider locale="id" messages={messages}>
      <LoginPage />
    </NextIntlClientProvider>,
  );
}

describe("LoginPage", () => {
  it("shows validation errors when submitting an empty form", async () => {
    const user = userEvent.setup();
    renderLoginPage();

    await user.click(screen.getByRole("button", { name: messages.Auth.loginCta }));

    expect(await screen.findByText(messages.Auth.emailInvalid)).toBeInTheDocument();
    expect(screen.getByText(messages.Auth.passwordMin)).toBeInTheDocument();
    expect(mutateMock).not.toHaveBeenCalled();
  });

  it("shows an email validation error for a malformed address", async () => {
    const user = userEvent.setup();
    renderLoginPage();

    await user.type(screen.getByLabelText(messages.Auth.email), "not-an-email");
    await user.type(screen.getByLabelText(messages.Auth.password), "Password123");
    await user.click(screen.getByRole("button", { name: messages.Auth.loginCta }));

    expect(await screen.findByText(messages.Auth.emailInvalid)).toBeInTheDocument();
    expect(mutateMock).not.toHaveBeenCalled();
  });

  it("submits valid credentials to the login mutation", async () => {
    const user = userEvent.setup();
    renderLoginPage();

    await user.type(screen.getByLabelText(messages.Auth.email), "user@example.com");
    await user.type(screen.getByLabelText(messages.Auth.password), "Password123");
    await user.click(screen.getByRole("button", { name: messages.Auth.loginCta }));

    await waitFor(() => {
      expect(mutateMock).toHaveBeenCalledWith(
        { email: "user@example.com", password: "Password123" },
        expect.anything(),
      );
    });
  });
});
