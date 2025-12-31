import { SetupProvider } from "./setup-context";
import { SetupWizard } from "./setup-wizard";

export default function Setup() {
  return (
    <SetupProvider>
      <SetupWizard />
    </SetupProvider>
  );
}
