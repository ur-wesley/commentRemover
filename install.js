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

function setupBinary() {
 const binDir = path.join(__dirname, "bin");

 if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
 }

 const binaryPath = path.join(binDir, BINARY_NAME + getBinaryExtension());

 try {
  console.log("🔨 Building from source...");
  const buildCmd = `go build -o "${binaryPath}" .`;
  execSync(buildCmd, {
   cwd: __dirname,
   stdio: "inherit",
   env: { ...process.env },
  });

  if (os.platform() !== "win32") {
   fs.chmodSync(binaryPath, "755");
  }

  console.log("✅ Successfully built commenter from source!");
  return;
 } catch (error) {
  console.log("⚠️  Go not found or build failed, trying prebuilt binary...");
 }

 const existingBinary = path.join(__dirname, "cr" + getBinaryExtension());
 if (fs.existsSync(existingBinary)) {
  console.log("📦 Using existing binary...");
  fs.copyFileSync(existingBinary, binaryPath);

  if (os.platform() !== "win32") {
   fs.chmodSync(binaryPath, "755");
  }

  console.log("✅ Successfully installed commenter!");
  return;
 }

 console.log("❌ No prebuilt binary found and Go build failed.");
 console.log("📝 Manual installation required:");
 console.log("   1. Install Go from https://golang.org/");
 console.log("   2. Run: go build -o commenter");
 console.log("   3. Move the binary to your PATH");
 process.exit(1);
}

function main() {
 console.log("🚀 Installing Comment Remover...");
 console.log(`Platform: ${getPlatform()}-${getArch()}`);

 try {
  setupBinary();

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
 } catch (error) {
  console.error("❌ Installation failed:", error.message);
  process.exit(1);
 }
}

if (require.main === module) {
 main();
}

module.exports = { main, getPlatform, getArch };
