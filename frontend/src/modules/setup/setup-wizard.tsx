import { Container, Stepper, Group, Button } from "@mantine/core";
import { useState, useRef } from "react";
import { useSetup } from "./setup-context";
import { AccountStep, type AccountStepRef } from "./account-step";
import { TokensStep, type TokensStepRef } from "./tokens-step";
import { AppsStep, type AppsStepRef } from "./apps-step";
import { CompletionStep } from "./completion-step";

export function SetupWizard() {
  const [active, setActive] = useState(0);
  const { updateSetupData } = useSetup();
  const accountStepRef = useRef<AccountStepRef>(null);
  const tokensStepRef = useRef<TokensStepRef>(null);
  const appsStepRef = useRef<AppsStepRef>(null);

  const nextStep = () => {
    if (active === 0 && accountStepRef.current) {
      if (accountStepRef.current.validate()) {
        setActive((current) => current + 1);
      }
    } else if (active === 1 && tokensStepRef.current) {
      if (tokensStepRef.current.validate()) {
        setActive((current) => current + 1);
      }
    } else if (active === 2 && appsStepRef.current) {
      if (appsStepRef.current.validate()) {
        setActive((current) => current + 1);
      }
    }
  };

  const prevStep = () =>
    setActive((current) => (current > 0 ? current - 1 : current));

  const handleAccountNext = (data: { username: string; password: string }) => {
    updateSetupData(data);
  };

  const handleTokensNext = (data: { githubToken: string }) => {
    updateSetupData(data);
  };

  const handleAppsNext = (data: { apps: { provider: string; key: string; command: string }[] }) => {
    updateSetupData(data);
  };

  return (
    <Container size="md" style={{ marginTop: 50 }}>
      <Stepper active={active} onStepClick={setActive}>
        <Stepper.Step label="Step 1" description="Account">
          <AccountStep ref={accountStepRef} onNext={handleAccountNext} />
        </Stepper.Step>
        <Stepper.Step label="Step 2" description="Tokens">
          <TokensStep ref={tokensStepRef} onNext={handleTokensNext} />
        </Stepper.Step>
        <Stepper.Step label="Step 3" description="Apps">
          <AppsStep ref={appsStepRef} onNext={handleAppsNext} />
        </Stepper.Step>
        <Stepper.Step label="Step 4" description="Complete">
          <CompletionStep />
        </Stepper.Step>
      </Stepper>

      {active < 3 && (
        <Group justify="center" mt="xl">
          <Button variant="default" onClick={prevStep} disabled={active === 0}>
            Back
          </Button>
          <Button onClick={nextStep}>
            {active === 2 ? "Finish" : "Next"}
          </Button>
        </Group>
      )}
    </Container>
  );
}
