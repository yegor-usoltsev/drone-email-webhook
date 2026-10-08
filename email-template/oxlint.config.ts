import { defineConfig } from "oxlint";
import core from "ultracite/oxlint/core";
import react from "ultracite/oxlint/react";

export default defineConfig({
  extends: [core, react],
  ignorePatterns: core.ignorePatterns,
  rules: {
    // The email reads top-down: the component first, then its parts and styles.
    "func-style": ["error", "declaration"],
    "no-use-before-define": ["error", { functions: false, variables: false }],
    "react/function-component-definition": [
      "error",
      { namedComponents: "function-declaration" },
    ],
  },
});
