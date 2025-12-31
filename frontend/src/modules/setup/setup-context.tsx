import { createContext, useContext, useState, type ReactNode } from "react";

interface SetupData {
  username: string;
  password: string;
}

interface SetupContextType {
  setupData: SetupData;
  updateSetupData: (data: Partial<SetupData>) => void;
}

const SetupContext = createContext<SetupContextType | undefined>(undefined);

export function SetupProvider({ children }: { children: ReactNode }) {
  const [setupData, setSetupData] = useState<SetupData>({
    username: "",
    password: "",
  });

  const updateSetupData = (data: Partial<SetupData>) => {
    setSetupData((prev) => ({ ...prev, ...data }));
  };

  return (
    <SetupContext.Provider value={{ setupData, updateSetupData }}>
      {children}
    </SetupContext.Provider>
  );
}

export function useSetup() {
  const context = useContext(SetupContext);
  if (context === undefined) {
    throw new Error("useSetup must be used within a SetupProvider");
  }
  return context;
}
