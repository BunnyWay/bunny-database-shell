#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const binaryPath = path.join(__dirname, `bunny-database-shell${ext}`);

try {
  execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
} catch (err) {
  if (err.status !== null) {
    process.exit(err.status);
  }
  console.error("Failed to run bunny-database-shell.");
  console.error("Run `npm rebuild bunny-database-shell` or reinstall the package.");
  process.exit(1);
}
