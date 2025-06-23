#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const https = require("https");
const { execSync } = require("child_process");
const os = require("os");

const REPO_URL = "https://github.com/ur-wesley/commentRemover";
const BINARY_NAME = "commenter";

function getPlatform() {
 const platform = os.platform();
 const arch = os.arch();

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

function getArch() {
 const arch = os.arch();
 switch (arch) {
  case "x64":
   return "amd64";
  case "arm64":
   return "arm64";
  case "ia32":
   return "386";
  default:
   return "amd64";
 }
}

function getBinaryExtension() {
 return os.platform() === "win32" ? ".exe" : "";
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

async function setupBinary() {
 const binDir = path.join(__dirname, "bin");

 if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
 }

 const binaryPath = path.join(binDir, BINARY_NAME + getBinaryExtension());

 // Check if binary already exists
 if (fs.existsSync(binaryPath)) {
  console.log("✅ Binary already exists, skipping download");
  return;
 }

 try {
  console.log("📦 Downloading prebuilt binary from GitHub releases...");

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

  console.log(
   "✅ Successfully downloaded and installed commenter from GitHub releases!"
  );
  return;
 } catch (error) {
  console.log(
   "⚠️  Failed to download from GitHub releases, checking for existing binary..."
  );

  // Check if there's an existing binary in the package (for development or manual placement)
  const possibleBinaries = [
   path.join(__dirname, BINARY_NAME + getBinaryExtension()),
   path.join(__dirname, "commenter" + getBinaryExtension()),
   path.join(__dirname, "commenter-darwin"),
   path.join(__dirname, "cr" + getBinaryExtension()),
  ];

  for (const existingBinary of possibleBinaries) {
   if (fs.existsSync(existingBinary)) {
    console.log("📦 Found existing binary, copying to bin directory...");
    fs.copyFileSync(existingBinary, binaryPath);

    if (os.platform() !== "win32") {
     fs.chmodSync(binaryPath, "755");
    }

    console.log("✅ Successfully installed commenter from existing binary!");
    return;
   }
  }

  console.log(
   "❌ Failed to download from GitHub releases and no existing binary found."
  );
  console.log("📝 Please ensure:");
  console.log("   1. You have internet connection");
  console.log(
   "   2. The GitHub release exists for version " + (await getLatestVersion())
  );
  console.log(
   "   3. Or manually download the binary from: " + REPO_URL + "/releases"
  );
  process.exit(1);
 }
}

async function main() {
 console.log("🚀 Installing Comment Remover from GitHub releases...");
 console.log(`Platform: ${getPlatform()}-${getArch()}`);

 try {
  await setupBinary();

  console.log("");
  console.log("🎉 Installation complete!");
  console.log("");
  console.log("Usage:");
  console.log("  commenter <file>           # Preview comment removal");
  console.log("  commenter -r src/          # Process directory recursively");
  console.log(
   "  commenter -w <file>        # Remove comments and save (short)"
  );
  console.log("  commenter -w -r project/   # Process and save recursively");
  console.log("  commenter -h               # Show help (short)");
  console.log("");
  console.log("Supported file types: .ts, .tsx, .js, .jsx, .go, .sql, .json");
  console.log("");
  console.log("💡 You can also run with: bun run commenter");
 } catch (error) {
  console.error("❌ Installation failed:", error.message);
  process.exit(1);
 }
}

if (require.main === module) {
 main();
}

module.exports = { main, getPlatform, getArch };
