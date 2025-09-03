import { defineConfig } from "eslint/config";
import typescriptEslint from "@typescript-eslint/eslint-plugin";
import globals from "globals";
import tsParser from "@typescript-eslint/parser";
import path from "node:path";
import { fileURLToPath } from "node:url";
import js from "@eslint/js";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all,
});

export default defineConfig([
  {
    extends: compat.extends(
      "eslint:recommended",
      "plugin:react/recommended",
      "plugin:@typescript-eslint/recommended",
      "plugin:prettier/recommended",
    ),

    plugins: {
      "@typescript-eslint": typescriptEslint,
    },

    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.jest,
      },

      parser: tsParser,
      ecmaVersion: 2021,
      sourceType: "module",

      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
    },

    settings: {
      react: {
        version: "detect",
      },
      "import/resolver": {
        alias: {
          map: [["@", "./src"]],
        },
      },
    },

    rules: {
      "prettier/prettier": "warn",
      "react/react-in-jsx-scope": "off",
      camelcase: "warn",
      "react/prop-types": "warn",
      "react/prefer-stateless-function": "warn",
      "class-methods-use-this": "off",
      "no-param-reassign": "warn",
      "no-plusplus": "warn",
      "react/jsx-props-no-spreading": "off",
      "react/static-property-placement": "warn",
      "prefer-destructuring": "warn",
      "react/forbid-prop-types": "warn",
      "no-use-before-define": "off",
      "react/no-array-index-key": "warn",
      "consistent-return": "warn",
      "react/require-default-props": "warn",
      "no-unused-expressions": "off",
      "no-underscore-dangle": "info",
      "@typescript-eslint/no-empty-function": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-unused-vars": "off",

      "no-empty": [
        2,
        {
          allowEmptyCatch: true,
        },
      ],

      "react/destructuring-assignment": "warn",
      "no-nested-ternary": "warn",
      "global-require": "off",
    },
  },
]);
