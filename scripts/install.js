const https = require("https");
const path = require("path");
const fs = require("fs");

const REPO = "BunnyWay/bunny-database-shell";

const PLATFORM_MAP = {
  "darwin-arm64": "bunny-database-shell_darwin_arm64",
  "darwin-x64": "bunny-database-shell_darwin_amd64",
  "linux-arm64": "bunny-database-shell_linux_arm64",
  "linux-x64": "bunny-database-shell_linux_amd64",
  "win32-x64": "bunny-database-shell_windows_amd64.exe",
};

const binDir = path.join(__dirname, "..", "bin");
const ext = process.platform === "win32" ? ".exe" : "";
const outputPath = path.join(binDir, `bunny-database-shell${ext}`);

if (fs.existsSync(outputPath)) {
  console.log("bunny-database-shell binary already exists, skipping download.");
  process.exit(0);
}

const key = `${process.platform}-${process.arch}`;
const binaryName = PLATFORM_MAP[key];

if (!binaryName) {
  console.error(`Unsupported platform: ${key}`);
  console.error(
    "You can build from source: go build -o bin/bunny-database-shell ./cmd/bunny-database-shell/"
  );
  process.exit(1);
}

const { version } = require(path.join(__dirname, "..", "package.json"));
const url = `https://github.com/${REPO}/releases/download/v${version}/${binaryName}`;

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

console.log(`Downloading bunny-database-shell v${version} for ${key}...`);

function download(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (
          res.statusCode >= 300 &&
          res.statusCode < 400 &&
          res.headers.location
        ) {
          return download(res.headers.location).then(resolve, reject);
        }
        if (res.statusCode !== 200) {
          return reject(new Error(`Download failed: HTTP ${res.statusCode}`));
        }
        const chunks = [];
        res.on("data", (chunk) => chunks.push(chunk));
        res.on("end", () => resolve(Buffer.concat(chunks)));
        res.on("error", reject);
      })
      .on("error", reject);
  });
}

download(url)
  .then((data) => {
    fs.writeFileSync(outputPath, data, { mode: 0o755 });
    console.log("bunny-database-shell installed successfully.");
  })
  .catch((err) => {
    console.error(`Failed to download binary: ${err.message}`);
    console.error(`URL: ${url}`);
    console.error(
      "You can build from source: go build -o bin/bunny-database-shell ./cmd/bunny-database-shell/"
    );
    process.exit(1);
  });
