import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { SetupProvider } from "./setup-context";
import { SetupWizard } from "./setup-wizard";

export default function Setup() {
  const [status, setStatus] = useState<"loading" | "pending-setup" | "other">("loading");

  useEffect(() => {
    fetch("/api/check")
      .then((res) => res.json())
      .then((data) => {
        if (data.status === "pending-setup") {
          setStatus("pending-setup");
        } else {
          setStatus("other");
        }
      })
      .catch(() => setStatus("other"));
  }, []);

  if (status === "loading") {
    return null;
  }

  if (status === "other") {
    return <Navigate to="/" replace />;
  }

  return (
    <SetupProvider>
      <SetupWizard />
    </SetupProvider>
  );
}
