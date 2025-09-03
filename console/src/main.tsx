import React from "react";
import ReactDOM from "react-dom/client";
import "react-windy-ui/dist/wui-dark_purple.css";
import "@/assets/styles/app.scss"
import "@/assets/i18n/i18n"
import App from "./App";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
    <App />
);
