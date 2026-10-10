import { defineConfig } from "oxlint";
import antiSlop from "ultracite/oxlint/anti-slop";
import core from "ultracite/oxlint/core";
import { jsPluginSettings, selectJsPlugins } from "ultracite/oxlint/js-plugins";
import react from "ultracite/oxlint/react";

const jsPlugins = selectJsPlugins(["sonarjs", "react-doctor", "github"]);

export default defineConfig({
  extends: [core, react, antiSlop, jsPlugins],
  ignorePatterns: [
    ...(core.ignorePatterns ?? []),
    ".agents/skills",
    ".claude/skills",
  ],
  jsPlugins: jsPlugins.jsPlugins,
  options: {
    reportUnusedDisableDirectives: "error",
    typeAware: true,
  },
  overrides: [
    {
      files: ["**/*.config.{cjs,cts,js,mjs,mts,ts}"],
      rules: {
        "import/no-default-export": "off",
      },
    },
    {
      files: ["**/*.tsx"],
      rules: {
        "func-style": ["error", "declaration"],
      },
    },
  ],
  rules: {
    "func-style": ["error", "expression"],
    "import/no-default-export": "error",
    // The email reads top-down: the component first, then its parts and styles.
    "no-use-before-define": ["error", { functions: false, variables: false }],
    "react/function-component-definition": [
      "error",
      {
        namedComponents: "function-declaration",
        unnamedComponents: "arrow-function",
      },
    ],
  },
  settings: jsPluginSettings,
});
