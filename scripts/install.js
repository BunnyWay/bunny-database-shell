const { execSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const binDir = path.join(__dirname, "..", "bin");
const outputPath = path.join(binDir, "bunny-database-shell");

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

const srcPath = path.join(__dirname, "..", "cmd", "bunny-database-shell");

console.log("Building bunny-database-shell...");

try {
  execSync(`go build -o ${outputPath} ${srcPath}`, {
    stdio: "inherit",
    cwd: path.join(__dirname, ".."),
  });
  console.log("bunny-database-shell installed successfully.");
} catch (err) {
  console.error("Failed to build bunny-database-shell.");
  console.error("Make sure Go is installed: https://go.dev/dl/");
  process.exit(1);
}
