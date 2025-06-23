#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");
const os = require("os");
const https = require("https");

const REPO_URL = "https://github.com/ur-wesley/commentRemover";

function getBinaryPath() {
 const platform = os.platform();
 const extension = platform === "win32" ? ".exe" : "";
 return path.join(__dirname, "bin", `commenter${extension}`);
}

function getPlatform() {
 const platform = os.platform();
 switch (platform) {
  case "win32":
   return "windows";
  case "darwin":
   return "darwin";
  case "linux":
   return "linux";
  default:
   throw new Error(`Unsupported platform: ${platform}`);
 }
}

function downloadBinary(url, dest) {
 return new Promise((resolve, reject) => {
  const file = fs.createWriteStream(dest);
  https
   .get(url, (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
     file.close();
     fs.unlink(dest, () => {});
     return downloadBinary(response.headers.location, dest)
      .then(resolve)
      .catch(reject);
    }
    if (response.statusCode !== 200) {
     file.close();
     fs.unlink(dest, () => {});
     reject(
      new Error(
       `HTTP ${response.statusCode}: ${response.statusMessage} for ${url}`
      )
     );
     return;
    }
    response.pipe(file);
    file.on("finish", () => {
     file.close();
     resolve();
    });
    file.on("error", (err) => {
     file.close();
     fs.unlink(dest, () => {});
     reject(err);
    });
   })
   .on("error", (err) => {
    file.close();
    fs.unlink(dest, () => {});
    reject(err);
   });
 });
}

async function getLatestVersion() {
 return new Promise((resolve, reject) => {
  const packageJsonPath = path.join(__dirname, "package.json");
  try {
   const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, "utf8"));
   resolve(packageJson.version);
  } catch (error) {
   reject(error);
  }
 });
}

async function ensureBinary() {
 const binaryPath = getBinaryPath();

 if (fs.existsSync(binaryPath)) {
  return binaryPath;
 }

 console.log("📦 Binary not found. Downloading from GitHub releases...");

 const binDir = path.dirname(binaryPath);
 if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
 }

 try {
  const version = await getLatestVersion();
  const platform = getPlatform();

  let binaryName;
  if (platform === "windows") {
   binaryName = "commenter.exe";
  } else if (platform === "darwin") {
   binaryName = "commenter-darwin";
  } else {
   binaryName = "commenter";
  }

  const downloadUrl = `${REPO_URL}/releases/download/v${version}/${binaryName}`;
  console.log(`Downloading from: ${downloadUrl}`);

  await downloadBinary(downloadUrl, binaryPath);

  if (os.platform() !== "win32") {
   fs.chmodSync(binaryPath, "755");
  }

  console.log("✅ Successfully downloaded commenter!");
  return binaryPath;
 } catch (error) {
  console.error("❌ Failed to download binary:", error.message);
  console.error(
   "Please ensure you have internet connection and the GitHub release exists."
  );
  process.exit(1);
 }
}

async function main() {
 try {
  const binaryPath = await ensureBinary();
  const args = process.argv.slice(2);

  const child = spawn(binaryPath, args, {
   stdio: "inherit",
   cwd: process.cwd(),
  });

  child.on("error", (error) => {
   console.error("❌ Failed to execute commenter:", error.message);
   process.exit(1);
  });

  child.on("close", (code) => {
   process.exit(code);
  });
 } catch (error) {
  console.error("❌ Error:", error.message);
  process.exit(1);
 }
}

if (require.main === module) {
 main();
}

module.exports = { main, getBinaryPath, ensureBinary };
